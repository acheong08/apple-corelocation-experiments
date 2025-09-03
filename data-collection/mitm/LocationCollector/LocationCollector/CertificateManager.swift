//
//  CertificateManager.swift
//  LocationCollector
//
//  Created by Antonio on 03/09/2025.
//

import Foundation
import Security
import os.log

class CertificateManager: ObservableObject {
    static let shared = CertificateManager()
    
    @Published var isRootCAGenerated = false
    @Published var certificateData: Data?
    
    private let logger = Logger(subsystem: "dev.duti.LocationCollector", category: "CertificateManager")
    private let keychain = Keychain()
    
    private let rootCALabel = "LocationCollector-RootCA"
    private let rootCACommonName = "LocationCollector Root CA"
    
    private init() {
        checkExistingCertificate()
    }
    
    private func checkExistingCertificate() {
        if let existingCert = keychain.getCertificate(label: rootCALabel) {
            logger.info("Found existing root CA certificate")
            certificateData = existingCert
            isRootCAGenerated = true
        } else {
            logger.info("No existing root CA found")
            isRootCAGenerated = false
        }
    }
    
    func generateRootCA() {
        logger.info("Generating root CA certificate")
        
        do {
            let (privateKey, certificate) = try createRootCACertificate()
            
            // Store in keychain
            try keychain.storePrivateKey(privateKey, label: "\(rootCALabel)-PrivateKey")
            try keychain.storeCertificate(certificate, label: rootCALabel)
            
            certificateData = certificate
            isRootCAGenerated = true
            
            logger.info("✅ Root CA certificate generated and stored successfully")
            
        } catch {
            logger.error("❌ Failed to generate root CA: \(error.localizedDescription)")
        }
    }
    
    private func createRootCACertificate() throws -> (SecKey, Data) {
        // Generate RSA key pair
        let keyAttributes: [String: Any] = [
            kSecAttrKeyType as String: kSecAttrKeyTypeRSA,
            kSecAttrKeySizeInBits as String: 2048,
            kSecPrivateKeyAttrs as String: [
                kSecAttrIsPermanent as String: false
            ]
        ]
        
        var error: Unmanaged<CFError>?
        guard let privateKey = SecKeyCreateRandomKey(keyAttributes as CFDictionary, &error) else {
            throw CertificateError.keyGenerationFailed(error?.takeRetainedValue())
        }
        
        guard let publicKey = SecKeyCopyPublicKey(privateKey) else {
            throw CertificateError.publicKeyExtractionFailed
        }
        
        // Create certificate
        let certificate = try createX509Certificate(privateKey: privateKey, publicKey: publicKey)
        
        return (privateKey, certificate)
    }
    
    private func createX509Certificate(privateKey: SecKey, publicKey: SecKey) throws -> Data {
        // For iOS, we'll use a simplified approach to create a self-signed certificate
        // This creates a basic DER-encoded certificate that iOS can recognize
        
        let now = Date()
        let oneYear = TimeInterval(365 * 24 * 60 * 60)
        let validFrom = now
        let validTo = now.addingTimeInterval(oneYear)
        
        // Get public key data
        guard let publicKeyData = SecKeyCopyExternalRepresentation(publicKey, nil) else {
            throw CertificateError.publicKeyDataExtractionFailed
        }
        
        // Create a minimal X.509 certificate using iOS Security framework
        let certificate = try createMinimalCertificate(
            privateKey: privateKey,
            publicKey: publicKey,
            publicKeyData: publicKeyData as Data,
            validFrom: validFrom,
            validTo: validTo,
            commonName: rootCACommonName
        )
        
        return certificate
    }
    
    private func createMinimalCertificate(privateKey: SecKey, publicKey: SecKey, publicKeyData: Data, validFrom: Date, validTo: Date, commonName: String) throws -> Data {
        // Create a basic certificate using SecCertificateCreateWithData approach
        // We'll build a minimal but valid X.509 certificate structure
        
        var tbsCertificate = Data()
        
        // Version (explicit v3)
        tbsCertificate.append(contentsOf: [0xa0, 0x03, 0x02, 0x01, 0x02])
        
        // Serial number (8 random bytes)
        let serialNumber = Data((0..<8).map { _ in UInt8.random(in: 0...255) })
        tbsCertificate.append(contentsOf: [0x02, 0x08])
        tbsCertificate.append(serialNumber)
        
        // Signature algorithm (SHA256withRSA)
        let sigAlgOID = Data([0x30, 0x0d, 0x06, 0x09, 0x2a, 0x86, 0x48, 0x86, 0xf7, 0x0d, 0x01, 0x01, 0x0b, 0x05, 0x00])
        tbsCertificate.append(sigAlgOID)
        
        // Issuer
        let issuer = createSimpleDN(cn: commonName)
        tbsCertificate.append(issuer)
        
        // Validity
        let validity = createSimpleValidity(from: validFrom, to: validTo)
        tbsCertificate.append(validity)
        
        // Subject (same as issuer for self-signed)
        tbsCertificate.append(issuer)
        
        // Subject Public Key Info
        let spki = createSimpleSPKI(publicKeyData: publicKeyData)
        tbsCertificate.append(spki)
        
        // Extensions (CA: TRUE)
        let extensions = createSimpleExtensions()
        tbsCertificate.append(extensions)
        
        // Wrap TBS certificate in SEQUENCE
        var tbsWrapper = Data()
        tbsWrapper.append(0x30) // SEQUENCE
        tbsWrapper.append(contentsOf: encodeLength(tbsCertificate.count))
        tbsWrapper.append(tbsCertificate)
        
        // Sign the TBS certificate
        let signature = try signData(tbsWrapper, with: privateKey)
        
        // Create final certificate
        var certificate = Data()
        certificate.append(0x30) // SEQUENCE
        
        let certLength = tbsWrapper.count + sigAlgOID.count + 4 + signature.count
        certificate.append(contentsOf: encodeLength(certLength))
        
        // TBS Certificate
        certificate.append(tbsWrapper)
        
        // Signature Algorithm
        certificate.append(sigAlgOID)
        
        // Signature
        certificate.append(0x03) // BIT STRING
        certificate.append(contentsOf: encodeLength(signature.count + 1))
        certificate.append(0x00) // No unused bits
        certificate.append(signature)
        
        return certificate
    }
    
    private func createSimpleDN(cn: String) -> Data {
        var dn = Data()
        let cnData = cn.data(using: .utf8)!
        
        dn.append(0x30) // SEQUENCE
        dn.append(contentsOf: encodeLength(cnData.count + 15))
        dn.append(0x31) // SET
        dn.append(contentsOf: encodeLength(cnData.count + 13))
        dn.append(0x30) // SEQUENCE
        dn.append(contentsOf: encodeLength(cnData.count + 11))
        dn.append(contentsOf: [0x06, 0x03, 0x55, 0x04, 0x03]) // CN OID (2.5.4.3)
        dn.append(0x0c) // UTF8String
        dn.append(contentsOf: encodeLength(cnData.count))
        dn.append(cnData)
        
        return dn
    }
    
    private func createSimpleValidity(from: Date, to: Date) -> Data {
        let formatter = DateFormatter()
        formatter.dateFormat = "yyMMddHHmmss'Z'"
        formatter.timeZone = TimeZone(identifier: "UTC")
        
        let notBefore = formatter.string(from: from).data(using: .ascii)!
        let notAfter = formatter.string(from: to).data(using: .ascii)!
        
        var validity = Data()
        validity.append(0x30) // SEQUENCE
        validity.append(contentsOf: encodeLength(notBefore.count + notAfter.count + 4))
        validity.append(0x17) // UTCTime
        validity.append(contentsOf: encodeLength(notBefore.count))
        validity.append(notBefore)
        validity.append(0x17) // UTCTime
        validity.append(contentsOf: encodeLength(notAfter.count))
        validity.append(notAfter)
        
        return validity
    }
    
    private func createSimpleSPKI(publicKeyData: Data) -> Data {
        var spki = Data()
        
        // Algorithm identifier for RSA encryption
        let rsaOID = Data([0x30, 0x0d, 0x06, 0x09, 0x2a, 0x86, 0x48, 0x86, 0xf7, 0x0d, 0x01, 0x01, 0x01, 0x05, 0x00])
        
        spki.append(0x30) // SEQUENCE
        let spkiContentLength = rsaOID.count + publicKeyData.count + 4
        spki.append(contentsOf: encodeLength(spkiContentLength))
        spki.append(rsaOID)
        spki.append(0x03) // BIT STRING
        spki.append(contentsOf: encodeLength(publicKeyData.count + 1))
        spki.append(0x00) // No unused bits
        spki.append(publicKeyData)
        
        return spki
    }
    
    private func createSimpleExtensions() -> Data {
        var extensions = Data()
        
        // Basic Constraints: CA:TRUE
        let basicConstraints = Data([0x30, 0x0f, 0x06, 0x03, 0x55, 0x1d, 0x13, 0x01, 0x01, 0xff, 0x04, 0x05, 0x30, 0x03, 0x01, 0x01, 0xff])
        
        // Key Usage: Key Cert Sign, CRL Sign
        let keyUsage = Data([0x30, 0x0e, 0x06, 0x03, 0x55, 0x1d, 0x0f, 0x01, 0x01, 0xff, 0x04, 0x04, 0x03, 0x02, 0x01, 0x06])
        
        extensions.append(0xa3) // Context-specific tag [3] for extensions
        let extensionsContentLength = basicConstraints.count + keyUsage.count + 2
        extensions.append(contentsOf: encodeLength(extensionsContentLength))
        extensions.append(0x30) // SEQUENCE
        extensions.append(contentsOf: encodeLength(basicConstraints.count + keyUsage.count))
        extensions.append(basicConstraints)
        extensions.append(keyUsage)
        
        return extensions
    }
    
    private func signData(_ data: Data, with privateKey: SecKey) throws -> Data {
        let algorithm = SecKeyAlgorithm.rsaSignatureMessagePKCS1v15SHA256
        
        var error: Unmanaged<CFError>?
        guard let signature = SecKeyCreateSignature(privateKey, algorithm, data as CFData, &error) else {
            throw CertificateError.signingFailed(error?.takeRetainedValue())
        }
        
        return signature as Data
    }
    
    // Helper function to encode ASN.1 lengths properly
    private func encodeLength(_ length: Int) -> Data {
        var encoded = Data()
        
        if length < 128 {
            // Short form
            encoded.append(UInt8(length))
        } else {
            // Long form
            if length < 256 {
                encoded.append(0x81) // Long form, 1 byte
                encoded.append(UInt8(length))
            } else if length < 65536 {
                encoded.append(0x82) // Long form, 2 bytes
                encoded.append(UInt8((length >> 8) & 0xff))
                encoded.append(UInt8(length & 0xff))
            } else {
                encoded.append(0x83) // Long form, 3 bytes
                encoded.append(UInt8((length >> 16) & 0xff))
                encoded.append(UInt8((length >> 8) & 0xff))
                encoded.append(UInt8(length & 0xff))
            }
        }
        
        return encoded
    }
    
    func getCertificateForExport() -> Data? {
        return certificateData
    }
    
    func exportCertificateAsPEM() -> String? {
        guard let certData = certificateData else { return nil }
        
        let base64Cert = certData.base64EncodedString(options: [.lineLength64Characters, .endLineWithLineFeed])
        return "-----BEGIN CERTIFICATE-----\n\(base64Cert)\n-----END CERTIFICATE-----"
    }
}

// MARK: - Keychain Helper

class Keychain {
    private let logger = Logger(subsystem: "dev.duti.LocationCollector", category: "Keychain")
    
    func storePrivateKey(_ privateKey: SecKey, label: String) throws {
        let query: [String: Any] = [
            kSecClass as String: kSecClassKey,
            kSecAttrApplicationTag as String: label.data(using: .utf8)!,
            kSecAttrKeyType as String: kSecAttrKeyTypeRSA,
            kSecAttrKeyClass as String: kSecAttrKeyClassPrivate,
            kSecValueRef as String: privateKey
        ]
        
        let status = SecItemAdd(query as CFDictionary, nil)
        if status != errSecSuccess && status != errSecDuplicateItem {
            throw CertificateError.keychainStorageFailed(status)
        }
    }
    
    func storeCertificate(_ certificate: Data, label: String) throws {
        let query: [String: Any] = [
            kSecClass as String: kSecClassCertificate,
            kSecAttrLabel as String: label,
            kSecValueData as String: certificate
        ]
        
        let status = SecItemAdd(query as CFDictionary, nil)
        if status != errSecSuccess && status != errSecDuplicateItem {
            throw CertificateError.keychainStorageFailed(status)
        }
    }
    
    func getCertificate(label: String) -> Data? {
        let query: [String: Any] = [
            kSecClass as String: kSecClassCertificate,
            kSecAttrLabel as String: label,
            kSecReturnData as String: true,
            kSecMatchLimit as String: kSecMatchLimitOne
        ]
        
        var result: AnyObject?
        let status = SecItemCopyMatching(query as CFDictionary, &result)
        
        if status == errSecSuccess {
            return result as? Data
        }
        
        return nil
    }
    
    func getPrivateKey(label: String) -> SecKey? {
        let query: [String: Any] = [
            kSecClass as String: kSecClassKey,
            kSecAttrApplicationTag as String: label.data(using: .utf8)!,
            kSecAttrKeyType as String: kSecAttrKeyTypeRSA,
            kSecAttrKeyClass as String: kSecAttrKeyClassPrivate,
            kSecReturnRef as String: true,
            kSecMatchLimit as String: kSecMatchLimitOne
        ]
        
        var result: AnyObject?
        let status = SecItemCopyMatching(query as CFDictionary, &result)
        
        if status == errSecSuccess {
            return result as! SecKey
        }
        
        return nil
    }
}

// MARK: - Errors

enum CertificateError: Error, LocalizedError {
    case keyGenerationFailed(CFError?)
    case publicKeyExtractionFailed
    case publicKeyDataExtractionFailed
    case signingFailed(CFError?)
    case keychainStorageFailed(OSStatus)
    
    var errorDescription: String? {
        switch self {
        case .keyGenerationFailed(let error):
            return "Key generation failed: \(error?.localizedDescription ?? "Unknown error")"
        case .publicKeyExtractionFailed:
            return "Failed to extract public key"
        case .publicKeyDataExtractionFailed:
            return "Failed to extract public key data"
        case .signingFailed(let error):
            return "Certificate signing failed: \(error?.localizedDescription ?? "Unknown error")"
        case .keychainStorageFailed(let status):
            return "Keychain storage failed with status: \(status)"
        }
    }
}
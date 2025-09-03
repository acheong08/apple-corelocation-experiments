//
//  CertificateServer.swift
//  LocationCollector
//
//  Created by Antonio on 03/09/2025.
//

import Foundation
import Network
import os.log

class CertificateServer: ObservableObject {
    static let shared = CertificateServer()
    
    @Published var isRunning = false
    @Published var serverURL: String = ""
    
    private var listener: NWListener?
    private let logger = Logger(subsystem: "dev.duti.LocationCollector", category: "CertificateServer")
    private let port: NWEndpoint.Port = 8080
    
    private init() {}
    
    func startServer() {
        guard !isRunning else { return }
        
        do {
            let parameters = NWParameters.tcp
            parameters.allowLocalEndpointReuse = true
            
            listener = try NWListener(using: parameters, on: port)
            
            listener?.newConnectionHandler = { [weak self] connection in
                self?.handleConnection(connection)
            }
            
            listener?.stateUpdateHandler = { [weak self] state in
                switch state {
                case .ready:
                    self?.logger.info("Certificate server started on port \(self?.port.rawValue ?? 0)")
                    DispatchQueue.main.async {
                        self?.isRunning = true
                        self?.serverURL = "http://10.0.0.1:\(self?.port.rawValue ?? 8080)"
                    }
                case .failed(let error):
                    self?.logger.error("Certificate server failed: \(error.localizedDescription)")
                    DispatchQueue.main.async {
                        self?.isRunning = false
                    }
                case .cancelled:
                    self?.logger.info("Certificate server cancelled")
                    DispatchQueue.main.async {
                        self?.isRunning = false
                    }
                default:
                    break
                }
            }
            
            listener?.start(queue: .global(qos: .userInitiated))
            
        } catch {
            logger.error("Failed to start certificate server: \(error.localizedDescription)")
        }
    }
    
    func stopServer() {
        listener?.cancel()
        listener = nil
        DispatchQueue.main.async {
            self.isRunning = false
            self.serverURL = ""
        }
    }
    
    private func handleConnection(_ connection: NWConnection) {
        connection.stateUpdateHandler = { state in
            switch state {
            case .ready:
                self.receiveRequest(connection)
            case .failed(let error):
                self.logger.error("Connection failed: \(error.localizedDescription)")
                connection.cancel()
            case .cancelled:
                break
            default:
                break
            }
        }
        
        connection.start(queue: .global(qos: .userInitiated))
    }
    
    private func receiveRequest(_ connection: NWConnection) {
        connection.receive(minimumIncompleteLength: 1, maximumLength: 8192) { [weak self] data, _, isComplete, error in
            if let error = error {
                self?.logger.error("Receive error: \(error.localizedDescription)")
                connection.cancel()
                return
            }
            
            guard let data = data, !data.isEmpty else {
                connection.cancel()
                return
            }
            
            if let request = String(data: data, encoding: .utf8) {
                self?.logger.info("Received request: \(request.prefix(100))")
                self?.handleHTTPRequest(request, connection: connection)
            }
            
            if isComplete {
                connection.cancel()
            }
        }
    }
    
    private func handleHTTPRequest(_ request: String, connection: NWConnection) {
        let lines = request.components(separatedBy: "\r\n")
        guard let requestLine = lines.first else {
            sendErrorResponse(connection, status: "400 Bad Request")
            return
        }
        
        let components = requestLine.components(separatedBy: " ")
        guard components.count >= 2 else {
            sendErrorResponse(connection, status: "400 Bad Request")
            return
        }
        
        let method = components[0]
        let path = components[1]
        
        logger.info("HTTP \(method) \(path)")
        
        switch (method, path) {
        case ("GET", "/"):
            sendLandingPage(connection)
        case ("GET", "/cert"):
            sendCertificate(connection)
        case ("GET", "/mitmproxy-ca-cert.pem"):
            sendCertificate(connection)  // Alternative endpoint name
        default:
            sendErrorResponse(connection, status: "404 Not Found")
        }
    }
    
    private func sendLandingPage(_ connection: NWConnection) {
        let html = """
        <!DOCTYPE html>
        <html>
        <head>
            <title>LocationCollector Certificate</title>
            <meta name="viewport" content="width=device-width, initial-scale=1">
            <style>
                body { font-family: -apple-system, BlinkMacSystemFont, sans-serif; margin: 40px; line-height: 1.6; }
                .container { max-width: 600px; margin: 0 auto; }
                .cert-button { 
                    display: inline-block; 
                    background: #007AFF; 
                    color: white; 
                    padding: 15px 30px; 
                    text-decoration: none; 
                    border-radius: 8px; 
                    font-weight: bold;
                    margin: 20px 0;
                }
                .steps { background: #f5f5f5; padding: 20px; border-radius: 8px; margin: 20px 0; }
                .step { margin: 10px 0; }
            </style>
        </head>
        <body>
            <div class="container">
                <h1>🔒 LocationCollector Certificate</h1>
                <p>To enable HTTPS interception for Apple location services, you need to install and trust our root certificate.</p>
                
                <a href="/cert" class="cert-button">📥 Download Certificate</a>
                
                <div class="steps">
                    <h3>Installation Steps:</h3>
                    <div class="step">1. Tap "Download Certificate" above</div>
                    <div class="step">2. Go to Settings → General → VPN & Device Management</div>
                    <div class="step">3. Install the "LocationCollector Root CA" profile</div>
                    <div class="step">4. Go to Settings → General → About → Certificate Trust Settings</div>
                    <div class="step">5. Enable full trust for "LocationCollector Root CA"</div>
                </div>
                
                <p><small>This certificate is only used for intercepting Apple location service requests on this device.</small></p>
            </div>
        </body>
        </html>
        """
        
        sendHTTPResponse(connection, status: "200 OK", contentType: "text/html", body: html.data(using: .utf8)!)
    }
    
    private func sendCertificate(_ connection: NWConnection) {
        guard let certManager = CertificateManager.shared.getCertificateForExport() else {
            sendErrorResponse(connection, status: "500 Internal Server Error", body: "Certificate not available")
            return
        }
        
        // Convert to PEM format if not already
        let pemCert: String
        if let pemString = CertificateManager.shared.exportCertificateAsPEM() {
            pemCert = pemString
        } else {
            sendErrorResponse(connection, status: "500 Internal Server Error", body: "Failed to export certificate")
            return
        }
        
        logger.info("Serving certificate (\(pemCert.count) characters)")
        
        sendHTTPResponse(
            connection,
            status: "200 OK",
            contentType: "application/x-pem-file",
            headers: [
                "Content-Disposition: attachment; filename=\"LocationCollector-CA.pem\""
            ],
            body: pemCert.data(using: .utf8)!
        )
    }
    
    private func sendErrorResponse(_ connection: NWConnection, status: String, body: String = "") {
        let html = """
        <!DOCTYPE html>
        <html>
        <head><title>\(status)</title></head>
        <body>
            <h1>\(status)</h1>
            <p>\(body.isEmpty ? status : body)</p>
        </body>
        </html>
        """
        
        sendHTTPResponse(connection, status: status, contentType: "text/html", body: html.data(using: .utf8)!)
    }
    
    private func sendHTTPResponse(_ connection: NWConnection, status: String, contentType: String, headers: [String] = [], body: Data) {
        var response = "HTTP/1.1 \(status)\r\n"
        response += "Content-Type: \(contentType)\r\n"
        response += "Content-Length: \(body.count)\r\n"
        response += "Connection: close\r\n"
        
        for header in headers {
            response += "\(header)\r\n"
        }
        
        response += "\r\n"
        
        var responseData = response.data(using: .utf8)!
        responseData.append(body)
        
        connection.send(content: responseData, completion: .contentProcessed { error in
            if let error = error {
                self.logger.error("Send error: \(error.localizedDescription)")
            }
            connection.cancel()
        })
    }
}
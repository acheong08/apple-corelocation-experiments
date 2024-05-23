
bool _CLPWifiAPLocationReadFrom(long self,long protobufString)

{
  byte bVar1;
  undefined8 uVar2;
  long *plVar3;
  long *plVar4;
  long *plVar5;
  byte bVar6;
  char cVar7;
  undefined4 uVar8;
  undefined8 uVar9;
  long serverHash;
  ulong uVar10;
  ulong uVar11;
  long *plVar12;
  ulong uVar13;
  undefined local_48 [16];
  long local_38;
  
  plVar5 = _DAT_7ffa41ed4678;
  plVar4 = _DAT_7ffa41ed4660;
  plVar3 = _DAT_7ffa41ed4648;
  plVar12 = _DAT_7ffa41ed4650;
  if (*(ulong *)(protobufString + *_DAT_7ffa41ed4678) <
      *(ulong *)(protobufString + *_DAT_7ffa41ed4660)) {
    while (uVar13 = 0, *(char *)(protobufString + *plVar12) == '\0') {
      bVar6 = 0;
      uVar11 = 0;
      do {
        serverHash = *(long *)(protobufString + *plVar5);
        uVar10 = serverHash + 1;
        if ((uVar10 == 0) || (*(ulong *)(protobufString + *plVar4) < uVar10)) {
          *(undefined *)(protobufString + *plVar12) = 1;
LAB_7ffa07293624:
          cVar7 = *(char *)(protobufString + *plVar12);
          if (cVar7 != '\0') {
            uVar11 = uVar13;
          }
          goto LAB_7ffa07293630;
        }
        bVar1 = *(byte *)(*(long *)(protobufString + *plVar3) + serverHash);
        *(ulong *)(protobufString + *plVar5) = uVar10;
        uVar11 = uVar11 | (ulong)(bVar1 & 0x7f) << (bVar6 & 0x3f);
        if (-1 < (char)bVar1) goto LAB_7ffa07293624;
        bVar6 = bVar6 + 7;
      } while (bVar6 != 0x46);
      cVar7 = *(char *)(protobufString + *plVar12);
      uVar11 = 0;
LAB_7ffa07293630:
      if ((cVar7 != '\0') || (bVar6 = (byte)uVar11 & 7, bVar6 == 4)) break;
      switch((int)(uVar11 >> 3)) {
      case 1:
        uVar9 = _PBReaderReadString(protobufString);
        uVar9 = _objc_retainAutoreleasedReturnValue(uVar9);
        serverHash = _OBJC_IVAR_$_CLPWifiAPLocation._mac;
        goto LAB_7ffa072937e9;
      case 2:
        bVar6 = 0;
        uVar11 = 0;
        do {
          serverHash = *(long *)(protobufString + *plVar5);
          uVar10 = serverHash + 1;
          if ((uVar10 == 0) || (*(ulong *)(protobufString + *plVar4) < uVar10)) {
            *(undefined *)(protobufString + *plVar12) = 1;
LAB_7ffa07293980:
            if (*(char *)(protobufString + *plVar12) != '\0') {
              uVar11 = uVar13;
            }
            goto LAB_7ffa0729398b;
          }
          bVar1 = *(byte *)(*(long *)(protobufString + *plVar3) + serverHash);
          *(ulong *)(protobufString + *plVar5) = uVar10;
          uVar11 = uVar11 | (ulong)(bVar1 & 0x7f) << (bVar6 & 0x3f);
          if (-1 < (char)bVar1) goto LAB_7ffa07293980;
          bVar6 = bVar6 + 7;
        } while (bVar6 != 0x46);
        uVar11 = 0;
LAB_7ffa0729398b:
        uVar8 = (undefined4)uVar11;
        serverHash = _OBJC_IVAR_$_CLPWifiAPLocation._channel;
        break;
      case 3:
        bVar6 = 0;
        uVar11 = 0;
        do {
          serverHash = *(long *)(protobufString + *plVar5);
          uVar10 = serverHash + 1;
          if ((uVar10 == 0) || (*(ulong *)(protobufString + *plVar4) < uVar10)) {
            *(undefined *)(protobufString + *plVar12) = 1;
LAB_7ffa0729399b:
            if (*(char *)(protobufString + *plVar12) != '\0') {
              uVar11 = uVar13;
            }
            goto LAB_7ffa072939a6;
          }
          bVar1 = *(byte *)(*(long *)(protobufString + *plVar3) + serverHash);
          *(ulong *)(protobufString + *plVar5) = uVar10;
          uVar11 = uVar11 | (ulong)(bVar1 & 0x7f) << (bVar6 & 0x3f);
          if (-1 < (char)bVar1) goto LAB_7ffa0729399b;
          bVar6 = bVar6 + 7;
        } while (bVar6 != 0x46);
        uVar11 = 0;
LAB_7ffa072939a6:
        uVar8 = (undefined4)uVar11;
        serverHash = _OBJC_IVAR_$_CLPWifiAPLocation._rssi;
        break;
      case 4:
        uVar9 = _objc_alloc_init(PTR__OBJC_CLASS_$_CLPLocation_7ffa405a9658);
        local_38 = self;
        _objc_storeStrong(_OBJC_IVAR_$_CLPWifiAPLocation._location + self,uVar9);
        cVar7 = _PBReaderPlaceMark(protobufString,local_48);
        if ((cVar7 == '\0') || (cVar7 = _CLPLocationReadFrom(uVar9,protobufString), cVar7 == '\0'))
        {
          (*_DAT_7ffa41ee65d0)(uVar9);
          return false;
        }
        _PBReaderRecallMark(protobufString);
        (*_DAT_7ffa41ee65d0)(uVar9);
        plVar12 = _DAT_7ffa41ed4650;
        self = local_38;
        goto LAB_7ffa07293a02;
      case 5:
        uVar9 = _PBReaderReadString(protobufString);
        uVar9 = _objc_retainAutoreleasedReturnValue(uVar9);
        serverHash = _OBJC_IVAR_$_CLPWifiAPLocation._appBundleId;
LAB_7ffa072937e9:
        uVar2 = *(undefined8 *)(self + serverHash);
        *(undefined8 *)(self + serverHash) = uVar9;
        (*_DAT_7ffa41ee65d0)(uVar2);
        plVar12 = _DAT_7ffa41ed4650;
        goto LAB_7ffa07293a02;
      case 6:
        *(byte *)(self + _OBJC_IVAR_$_CLPWifiAPLocation._has) =
             *(byte *)(self + _OBJC_IVAR_$_CLPWifiAPLocation._has) | 8;
        bVar6 = 0;
        uVar11 = 0;
        do {
          serverHash = *(long *)(protobufString + *plVar5);
          uVar10 = serverHash + 1;
          if ((uVar10 == 0) || (*(ulong *)(protobufString + *plVar4) < uVar10)) {
            *(undefined *)(protobufString + *plVar12) = 1;
LAB_7ffa072939b6:
            if (*(char *)(protobufString + *plVar12) != '\0') {
              uVar11 = uVar13;
            }
            goto LAB_7ffa072939c1;
          }
          bVar1 = *(byte *)(*(long *)(protobufString + *plVar3) + serverHash);
          *(ulong *)(protobufString + *plVar5) = uVar10;
          uVar11 = uVar11 | (ulong)(bVar1 & 0x7f) << (bVar6 & 0x3f);
          if (-1 < (char)bVar1) goto LAB_7ffa072939b6;
          bVar6 = bVar6 + 7;
        } while (bVar6 != 0x46);
        uVar11 = 0;
LAB_7ffa072939c1:
        uVar8 = (undefined4)uVar11;
        serverHash = _OBJC_IVAR_$_CLPWifiAPLocation._serverHash;
        break;
      case 7:
        *(byte *)(self + _OBJC_IVAR_$_CLPWifiAPLocation._has) =
             *(byte *)(self + _OBJC_IVAR_$_CLPWifiAPLocation._has) | 2;
        bVar6 = 0;
        uVar11 = 0;
        do {
          serverHash = *(long *)(protobufString + *plVar5);
          uVar10 = serverHash + 1;
          if ((uVar10 == 0) || (*(ulong *)(protobufString + *plVar4) < uVar10)) {
            *(undefined *)(protobufString + *plVar12) = 1;
LAB_7ffa072939d1:
            if (*(char *)(protobufString + *plVar12) != '\0') {
              uVar11 = uVar13;
            }
            goto LAB_7ffa072939dc;
          }
          bVar1 = *(byte *)(*(long *)(protobufString + *plVar3) + serverHash);
          *(ulong *)(protobufString + *plVar5) = uVar10;
          uVar11 = uVar11 | (ulong)(bVar1 & 0x7f) << (bVar6 & 0x3f);
          if (-1 < (char)bVar1) goto LAB_7ffa072939d1;
          bVar6 = bVar6 + 7;
        } while (bVar6 != 0x46);
        uVar11 = 0;
LAB_7ffa072939dc:
        uVar8 = (undefined4)uVar11;
        serverHash = _OBJC_IVAR_$_CLPWifiAPLocation._hidden;
        break;
      case 8:
        *(byte *)(self + _OBJC_IVAR_$_CLPWifiAPLocation._has) =
             *(byte *)(self + _OBJC_IVAR_$_CLPWifiAPLocation._has) | 1;
        uVar13 = *(ulong *)(protobufString + *plVar5);
        if ((uVar13 < 0xfffffffffffffff8) && (uVar13 + 8 <= *(ulong *)(protobufString + *plVar4))) {
          uVar9 = *(undefined8 *)(*(long *)(protobufString + *plVar3) + uVar13);
          *(ulong *)(protobufString + *plVar5) = uVar13 + 8;
        }
        else {
          *(undefined *)(protobufString + *plVar12) = 1;
          uVar9 = 0;
        }
        *(undefined8 *)(self + _OBJC_IVAR_$_CLPWifiAPLocation._scanTimestamp) = uVar9;
        goto LAB_7ffa07293a02;
      case 9:
        *(byte *)(self + _OBJC_IVAR_$_CLPWifiAPLocation._has) =
             *(byte *)(self + _OBJC_IVAR_$_CLPWifiAPLocation._has) | 4;
        bVar6 = 0;
        uVar11 = 0;
        do {
          serverHash = *(long *)(protobufString + *plVar5);
          uVar10 = serverHash + 1;
          if ((uVar10 == 0) || (*(ulong *)(protobufString + *plVar4) < uVar10)) {
            *(undefined *)(protobufString + *plVar12) = 1;
LAB_7ffa072939ec:
            if (*(char *)(protobufString + *plVar12) != '\0') {
              uVar11 = uVar13;
            }
            goto LAB_7ffa072939f7;
          }
          bVar1 = *(byte *)(*(long *)(protobufString + *plVar3) + serverHash);
          *(ulong *)(protobufString + *plVar5) = uVar10;
          uVar11 = uVar11 | (ulong)(bVar1 & 0x7f) << (bVar6 & 0x3f);
          if (-1 < (char)bVar1) goto LAB_7ffa072939ec;
          bVar6 = bVar6 + 7;
        } while (bVar6 != 0x46);
        uVar11 = 0;
LAB_7ffa072939f7:
        uVar8 = (undefined4)uVar11;
        serverHash = _OBJC_IVAR_$_CLPWifiAPLocation._scanType;
        break;
      default:
        cVar7 = _PBReaderSkipValueWithTag(protobufString,uVar11 >> 3,bVar6);
        plVar12 = _DAT_7ffa41ed4650;
        if (cVar7 == '\0') {
          return false;
        }
        goto LAB_7ffa07293a02;
      }
      *(undefined4 *)(self + serverHash) = uVar8;
LAB_7ffa07293a02:
      if (*(ulong *)(protobufString + *plVar4) <= *(ulong *)(protobufString + *plVar5)) break;
    }
  }
  return *(char *)(protobufString + *_DAT_7ffa41ed4650) == '\0';
}

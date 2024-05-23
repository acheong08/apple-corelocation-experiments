
/* WARNING: Globals starting with '_' overlap smaller symbols at the same address */

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
  long tmpVar;
  ulong uVar10;
  ulong fieldType;
  long *plVar11;
  ulong uVar12;
  undefined local_48 [16];
  long local_38;
  
  plVar5 = _DAT_7ffa41ed4678;
  plVar4 = _DAT_7ffa41ed4660;
  plVar3 = _DAT_7ffa41ed4648;
  plVar11 = _DAT_7ffa41ed4650;
  if (*(ulong *)(protobufString + *_DAT_7ffa41ed4678) <
      *(ulong *)(protobufString + *_DAT_7ffa41ed4660)) {
    while (uVar12 = 0, *(char *)(protobufString + *plVar11) == '\0') {
      bVar6 = 0;
      fieldType = 0;
      do {
        tmpVar = *(long *)(protobufString + *plVar5);
        uVar10 = tmpVar + 1;
        if ((uVar10 == 0) || (*(ulong *)(protobufString + *plVar4) < uVar10)) {
          *(undefined *)(protobufString + *plVar11) = 1;
LAB_7ffa07293624:
          cVar7 = *(char *)(protobufString + *plVar11);
          if (cVar7 != '\0') {
            fieldType = uVar12;
          }
          goto parsingLoopStart;
        }
        bVar1 = *(byte *)(*(long *)(protobufString + *plVar3) + tmpVar);
        *(ulong *)(protobufString + *plVar5) = uVar10;
        fieldType = fieldType | (ulong)(bVar1 & 0x7f) << (bVar6 & 0x3f);
        if (-1 < (char)bVar1) goto LAB_7ffa07293624;
        bVar6 = bVar6 + 7;
      } while (bVar6 != 0x46);
      cVar7 = *(char *)(protobufString + *plVar11);
      fieldType = 0;
parsingLoopStart:
      if ((cVar7 != '\0') || (bVar6 = (byte)fieldType & 7, bVar6 == 4)) break;
      switch((int)(fieldType >> 3)) {
      case 1:
        uVar9 = _PBReaderReadString(protobufString);
        uVar9 = _objc_retainAutoreleasedReturnValue(uVar9);
        tmpVar = _OBJC_IVAR_$_CLPWifiAPLocation._mac;
        goto LAB_7ffa072937e9;
      case 2:
        bVar6 = 0;
        fieldType = 0;
        do {
          tmpVar = *(long *)(protobufString + *plVar5);
          uVar10 = tmpVar + 1;
          if ((uVar10 == 0) || (*(ulong *)(protobufString + *plVar4) < uVar10)) {
            *(undefined *)(protobufString + *plVar11) = 1;
LAB_7ffa07293980:
            if (*(char *)(protobufString + *plVar11) != '\0') {
              fieldType = uVar12;
            }
            goto LAB_7ffa0729398b;
          }
          bVar1 = *(byte *)(*(long *)(protobufString + *plVar3) + tmpVar);
          *(ulong *)(protobufString + *plVar5) = uVar10;
          fieldType = fieldType | (ulong)(bVar1 & 0x7f) << (bVar6 & 0x3f);
          if (-1 < (char)bVar1) goto LAB_7ffa07293980;
          bVar6 = bVar6 + 7;
        } while (bVar6 != 0x46);
        fieldType = 0;
LAB_7ffa0729398b:
        uVar8 = (undefined4)fieldType;
        tmpVar = _OBJC_IVAR_$_CLPWifiAPLocation._channel;
        break;
      case 3:
        bVar6 = 0;
        fieldType = 0;
        do {
          tmpVar = *(long *)(protobufString + *plVar5);
          uVar10 = tmpVar + 1;
          if ((uVar10 == 0) || (*(ulong *)(protobufString + *plVar4) < uVar10)) {
            *(undefined *)(protobufString + *plVar11) = 1;
LAB_7ffa0729399b:
            if (*(char *)(protobufString + *plVar11) != '\0') {
              fieldType = uVar12;
            }
            goto LAB_7ffa072939a6;
          }
          bVar1 = *(byte *)(*(long *)(protobufString + *plVar3) + tmpVar);
          *(ulong *)(protobufString + *plVar5) = uVar10;
          fieldType = fieldType | (ulong)(bVar1 & 0x7f) << (bVar6 & 0x3f);
          if (-1 < (char)bVar1) goto LAB_7ffa0729399b;
          bVar6 = bVar6 + 7;
        } while (bVar6 != 0x46);
        fieldType = 0;
LAB_7ffa072939a6:
        uVar8 = (undefined4)fieldType;
        tmpVar = _OBJC_IVAR_$_CLPWifiAPLocation._rssi;
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
        plVar11 = _DAT_7ffa41ed4650;
        self = local_38;
        goto LAB_7ffa07293a02;
      case 5:
        uVar9 = _PBReaderReadString(protobufString);
        uVar9 = _objc_retainAutoreleasedReturnValue(uVar9);
        tmpVar = _OBJC_IVAR_$_CLPWifiAPLocation._appBundleId;
LAB_7ffa072937e9:
        uVar2 = *(undefined8 *)(self + tmpVar);
        *(undefined8 *)(self + tmpVar) = uVar9;
        (*_DAT_7ffa41ee65d0)(uVar2);
        plVar11 = _DAT_7ffa41ed4650;
        goto LAB_7ffa07293a02;
      case 6:
        *(byte *)(self + _OBJC_IVAR_$_CLPWifiAPLocation._has) =
             *(byte *)(self + _OBJC_IVAR_$_CLPWifiAPLocation._has) | 8;
        bVar6 = 0;
        fieldType = 0;
        do {
          tmpVar = *(long *)(protobufString + *plVar5);
          uVar10 = tmpVar + 1;
          if ((uVar10 == 0) || (*(ulong *)(protobufString + *plVar4) < uVar10)) {
            *(undefined *)(protobufString + *plVar11) = 1;
LAB_7ffa072939b6:
            if (*(char *)(protobufString + *plVar11) != '\0') {
              fieldType = uVar12;
            }
            goto LAB_7ffa072939c1;
          }
          bVar1 = *(byte *)(*(long *)(protobufString + *plVar3) + tmpVar);
          *(ulong *)(protobufString + *plVar5) = uVar10;
          fieldType = fieldType | (ulong)(bVar1 & 0x7f) << (bVar6 & 0x3f);
          if (-1 < (char)bVar1) goto LAB_7ffa072939b6;
          bVar6 = bVar6 + 7;
        } while (bVar6 != 0x46);
        fieldType = 0;
LAB_7ffa072939c1:
        uVar8 = (undefined4)fieldType;
        tmpVar = _OBJC_IVAR_$_CLPWifiAPLocation._serverHash;
        break;
      case 7:
        *(byte *)(self + _OBJC_IVAR_$_CLPWifiAPLocation._has) =
             *(byte *)(self + _OBJC_IVAR_$_CLPWifiAPLocation._has) | 2;
        bVar6 = 0;
        fieldType = 0;
        do {
          tmpVar = *(long *)(protobufString + *plVar5);
          uVar10 = tmpVar + 1;
          if ((uVar10 == 0) || (*(ulong *)(protobufString + *plVar4) < uVar10)) {
            *(undefined *)(protobufString + *plVar11) = 1;
LAB_7ffa072939d1:
            if (*(char *)(protobufString + *plVar11) != '\0') {
              fieldType = uVar12;
            }
            goto LAB_7ffa072939dc;
          }
          bVar1 = *(byte *)(*(long *)(protobufString + *plVar3) + tmpVar);
          *(ulong *)(protobufString + *plVar5) = uVar10;
          fieldType = fieldType | (ulong)(bVar1 & 0x7f) << (bVar6 & 0x3f);
          if (-1 < (char)bVar1) goto LAB_7ffa072939d1;
          bVar6 = bVar6 + 7;
        } while (bVar6 != 0x46);
        fieldType = 0;
LAB_7ffa072939dc:
        uVar8 = (undefined4)fieldType;
        tmpVar = _OBJC_IVAR_$_CLPWifiAPLocation._hidden;
        break;
      case 8:
        *(byte *)(self + _OBJC_IVAR_$_CLPWifiAPLocation._has) =
             *(byte *)(self + _OBJC_IVAR_$_CLPWifiAPLocation._has) | 1;
        uVar12 = *(ulong *)(protobufString + *plVar5);
        if ((uVar12 < 0xfffffffffffffff8) && (uVar12 + 8 <= *(ulong *)(protobufString + *plVar4))) {
          uVar9 = *(undefined8 *)(*(long *)(protobufString + *plVar3) + uVar12);
          *(ulong *)(protobufString + *plVar5) = uVar12 + 8;
        }
        else {
          *(undefined *)(protobufString + *plVar11) = 1;
          uVar9 = 0;
        }
        *(undefined8 *)(self + _OBJC_IVAR_$_CLPWifiAPLocation._scanTimestamp) = uVar9;
        goto LAB_7ffa07293a02;
      case 9:
        *(byte *)(self + _OBJC_IVAR_$_CLPWifiAPLocation._has) =
             *(byte *)(self + _OBJC_IVAR_$_CLPWifiAPLocation._has) | 4;
        bVar6 = 0;
        fieldType = 0;
        do {
          tmpVar = *(long *)(protobufString + *plVar5);
          uVar10 = tmpVar + 1;
          if ((uVar10 == 0) || (*(ulong *)(protobufString + *plVar4) < uVar10)) {
            *(undefined *)(protobufString + *plVar11) = 1;
LAB_7ffa072939ec:
            if (*(char *)(protobufString + *plVar11) != '\0') {
              fieldType = uVar12;
            }
            goto LAB_7ffa072939f7;
          }
          bVar1 = *(byte *)(*(long *)(protobufString + *plVar3) + tmpVar);
          *(ulong *)(protobufString + *plVar5) = uVar10;
          fieldType = fieldType | (ulong)(bVar1 & 0x7f) << (bVar6 & 0x3f);
          if (-1 < (char)bVar1) goto LAB_7ffa072939ec;
          bVar6 = bVar6 + 7;
        } while (bVar6 != 0x46);
        fieldType = 0;
LAB_7ffa072939f7:
        uVar8 = (undefined4)fieldType;
        tmpVar = _OBJC_IVAR_$_CLPWifiAPLocation._scanType;
        break;
      default:
        cVar7 = _PBReaderSkipValueWithTag(protobufString,fieldType >> 3,bVar6);
        plVar11 = _DAT_7ffa41ed4650;
        if (cVar7 == '\0') {
          return false;
        }
        goto LAB_7ffa07293a02;
      }
      *(undefined4 *)(self + tmpVar) = uVar8;
LAB_7ffa07293a02:
      if (*(ulong *)(protobufString + *plVar4) <= *(ulong *)(protobufString + *plVar5)) break;
    }
  }
  return *(char *)(protobufString + *_DAT_7ffa41ed4650) == '\0';
}


void _GEOTileKeyMake(byte *tileKeyPointer, ulong param_2, uint param_3,
                     uint param_4, uint param_5, byte param_6, uint param_7,
                     byte param_8, char param_9)

{
  *(ulong *)(tileKeyPointer + 1) =
      (ulong)(param_4 & 0x3f) << 0x28 | param_2 << 0x2e;
  *(uint *)(tileKeyPointer + 9) =
      (param_3 & 0x3ffffff) << 8 | (uint)(param_2 >> 0x12) & 0xff;
  tileKeyPointer[0xf] =
      (byte)(((ulong)param_7 << 0x34) >> 0x30) | param_6 & 0xf;
  *(ushort *)(tileKeyPointer + 0xd) =
      (ushort)(((ulong)(param_5 & 0x3fff) << 0x22) >> 0x20) |
      (ushort)(byte)((param_3 & 0x3ffffff) >> 0x18);
  *tileKeyPointer = param_9 << 7 | param_8 & 0x7f;
  return;
}

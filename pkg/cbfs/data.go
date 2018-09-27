package cbfs

var (
	Master = []byte{
		// Put some data in to make sure we don't get fooled by starting at 0.
		/*00000000*/ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		/*00000010*/ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		/*00000020*/ 0x4c, 0x41, 0x52, 0x43, 0x48, 0x49, 0x56, 0x45, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x02, //|LARCHIVE... ....|
		/*00000030*/ 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x38, 0x63, 0x62, 0x66, 0x73, 0x20, 0x6d, 0x61, 0x73, //|.......8cbfs mas|
		/*00000040*/ 0x74, 0x65, 0x72, 0x20, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, //|ter header......|
		/*00000050*/ 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x4f, 0x52, 0x42, 0x43, 0x31, 0x31, 0x31, 0x32, //|........ORBC1112|
		/*00000060*/ 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x02, 0x00, //|...........@....|
		/*00000070*/ 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, //|................|
		/*00000080*/ 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, //|................|
		/*00000090*/ 0x4c, 0x41, 0x52, 0x43, 0x48, 0x49, 0x56, 0x45, 0x00, 0x00, 0x00, 0x20, 0xff, 0xff, 0xff, 0xff, //|LARCHIVE..w.....|
		/*000000a0*/ 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x28, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, //|.......(........|
		/*000000b0*/ 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, //|................|
		/*000000c0*/ 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, //|................|
		/*000000d0*/ 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, //|................|

	}
	ListOutput = `FMAP REGION: COREBOOT
Name                           Offset     Type           Size   Comp
cbfs master header             0x0        cbfs header        32 none
fallback/romstage              0x80       stage           15812 none
fallback/ramstage              0x3ec0     stage           52417 none
config                         0x10bc0    raw               355 none
revision                       0x10d80    raw               576 none
cmos_layout.bin                0x11000    cmos_layout       548 none
fallback/dsdt.aml              0x11280    raw              6952 none
fallback/payload               0x12e00    simple elf         28 none
(empty)                        0x12e80    null           183192 none
bootblock                      0x3fa40    bootblock         880 none
`
)

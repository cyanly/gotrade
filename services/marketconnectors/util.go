package marketconnectors

var protoEnumToFIXEnum = map[int]string{
	0:  "0",
	1:  "1",
	2:  "2",
	3:  "3",
	4:  "4",
	5:  "5",
	6:  "6",
	7:  "7",
	8:  "8",
	9:  "9",
	10: "A",
	11: "B",
	12: "C",
	13: "D",
	14: "E",
	15: "F",
	16: "G",
	17: "H",
	18: "I",
	19: "J",
	20: "K",
	21: "L",
	22: "M",
	23: "N",
	24: "O",
	25: "P",
	26: "Q",
}
var fixEnumToProtoEnum = map[string]int{}

func init() {
	for k, v := range protoEnumToFIXEnum {
		fixEnumToProtoEnum[v] = k
	}
}

func ProtoEnumToFIXEnum(pbe int) string {
	return protoEnumToFIXEnum[pbe]
}
func FIXEnumToProtoEnum(fbe string) int {
	return fixEnumToProtoEnum[fbe]
}

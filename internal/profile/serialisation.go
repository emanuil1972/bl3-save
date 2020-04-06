package profile

import (
	"bufio"
	"io"

	"github.com/cfi2017/bl3-save/internal/pb"
	"github.com/cfi2017/bl3-save/internal/shared"
	"google.golang.org/protobuf/proto"
)

var (
	prefixMagic = []byte{
		0xD8, 0x04, 0xB9, 0x08, 0x5C, 0x4E, 0x2B, 0xC0,
		0x61, 0x9F, 0x7C, 0x8D, 0x5D, 0x34, 0x00, 0x56,
		0xE7, 0x7B, 0x4E, 0xC0, 0xA4, 0xD6, 0xA7, 0x01,
		0x14, 0x15, 0xA9, 0x93, 0x1F, 0x27, 0x2C, 0x8F,
	}
	xorMagic = []byte{
		0xE8, 0xDC, 0x3A, 0x66, 0xF7, 0xEF, 0x85, 0xE0,
		0xBD, 0x4A, 0xA9, 0x73, 0x57, 0x99, 0x30, 0x8C,
		0x94, 0x63, 0x59, 0xA8, 0xC9, 0xAE, 0xD9, 0x58,
		0x7D, 0x51, 0xB0, 0x1E, 0xBE, 0xD0, 0x77, 0x43,
	}
)

func Deserialize(reader io.Reader) (shared.SavFile, pb.Profile) {

	// deserialise header, decrypt data
	s, data := shared.DeserializeHeader(reader)

	data = shared.Decrypt(data, prefixMagic, xorMagic)

	p := pb.Profile{}
	if err := proto.Unmarshal(data, &p); err != nil {
		panic("couldn't unmarshal protobuf data")
	}

	return s, p
}

func Serialize(writer io.Writer, s shared.SavFile, p pb.Profile) {
	w := bufio.NewWriter(writer)
	shared.WriteBytes(w, []byte("GVAS"))
	shared.WriteInt(w, s.SgVersion)
	shared.WriteInt(w, s.PkgVersion)
	shared.WriteShort(w, s.EngineMajorVersion)
	shared.WriteShort(w, s.EngineMinorVersion)
	shared.WriteShort(w, s.EnginePatchVersion)
	shared.WriteInt(w, s.EngineBuildVersion)
	shared.WriteString(w, s.BuildId)
	shared.WriteInt(w, s.FmtVersion)
	shared.WriteInt(w, len(s.CustomFmtData))
	for _, d := range s.CustomFmtData {
		shared.WriteGuid(w, d.Guid)
		shared.WriteInt(w, d.Entry)
	}
	shared.WriteString(w, s.SgType)

	bs, err := proto.Marshal(&p)
	if err != nil {
		panic(err)
	}

	bs = shared.Encrypt(bs, prefixMagic, xorMagic)

	shared.WriteInt(w, len(bs))
	shared.WriteBytes(w, bs)

}

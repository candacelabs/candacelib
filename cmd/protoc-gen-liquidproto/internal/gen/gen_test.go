package gen

import (
	"go/parser"
	"go/token"
	"testing"

	liquidv1 "github.com/candacelabs/candacelib/liquidproto/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func TestGenerator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Liquid Proto generator")
}

var _ = Describe("Liquid Proto generator", func() {
	It("emits a direct message validator", func() {
		options := new(descriptorpb.FieldOptions)
		proto.SetExtension(options, liquidv1.E_Field, &liquidv1.FieldRefinement{
			Expr: "len(this) == 32",
		})
		fixture := &descriptorpb.FileDescriptorProto{
			Name:       proto.String("liquid_bytes_fixture.proto"),
			Package:    proto.String("candace.liquid.fixture.v1"),
			Syntax:     proto.String("proto3"),
			Dependency: []string{"liquidproto/v1/refinement.proto"},
			Options: &descriptorpb.FileOptions{
				GoPackage: proto.String("example.com/liquidfixture;liquidfixture"),
			},
			MessageType: []*descriptorpb.DescriptorProto{{
				Name: proto.String("Blob"),
				Field: []*descriptorpb.FieldDescriptorProto{{
					Name:    proto.String("digest"),
					Number:  proto.Int32(1),
					Label:   descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
					Type:    descriptorpb.FieldDescriptorProto_TYPE_BYTES.Enum(),
					Options: options,
				}},
			}},
		}
		request := &pluginpb.CodeGeneratorRequest{
			FileToGenerate: []string{fixture.GetName()},
			ProtoFile: []*descriptorpb.FileDescriptorProto{
				protodesc.ToFileDescriptorProto(descriptorpb.File_google_protobuf_descriptor_proto),
				protodesc.ToFileDescriptorProto(liquidv1.File_liquidproto_v1_refinement_proto),
				fixture,
			},
		}
		plugin, err := (protogen.Options{}).New(request)
		Expect(err).NotTo(HaveOccurred())
		Expect(Run(plugin)).To(Succeed())

		files := plugin.Response().GetFile()
		Expect(files).To(HaveLen(1))
		generated := files[0].GetContent()
		for _, want := range []string{
			"func ValidateBlob(message *Blob) error",
			"len(message.Digest) == 32",
		} {
			Expect(generated).To(ContainSubstring(want))
		}
		Expect(generated).To(MatchRegexp(`Value:\s+message[.]Digest`))
		for _, unwanted := range []string{"RefinedBlob", "MustBlob", "ToProto"} {
			Expect(generated).NotTo(ContainSubstring(unwanted))
		}

		_, err = parser.ParseFile(token.NewFileSet(), files[0].GetName(), generated, parser.AllErrors)
		Expect(err).NotTo(HaveOccurred(), "generated source is not valid Go syntax:\n%s", generated)
	})
})

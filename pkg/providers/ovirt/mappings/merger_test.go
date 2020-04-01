package mappings_test

import (
	v2vv1alpha1 "github.com/kubevirt/vm-import-operator/pkg/apis/v2v/v1alpha1"
	"github.com/kubevirt/vm-import-operator/pkg/providers/ovirt/mappings"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var (
	id1 = "id1"
	id2 = "id2"
	id3 = "id3"
	id4 = "id4"

	name1 = "name1"
	name2 = "name2"
	name3 = "name3"
	name4 = "name4"

	type1 = "type1"
	type2 = "type2"
)
var _ = Describe("Mappings merging ", func() {
	It("Should merge no mappings", func() {
		result := mappings.MergeMappings(nil, nil)

		Expect(result).To(Not(BeNil()))
		Expect(result.NetworkMappings).To(BeNil())
		Expect(result.StorageMappings).To(BeNil())
	})
	It("should produce nil mapping itemst on both input mapping items nil", func() {
		mapping := v2vv1alpha1.OvirtMappings{
			NetworkMappings: nil,
			StorageMappings: nil,
		}

		externalMapping := v2vv1alpha1.OvirtMappings{
			NetworkMappings: nil,
			StorageMappings: nil,
		}

		spec := v2vv1alpha1.ResourceMappingSpec{
			OvirtMappings: &externalMapping,
		}

		result := mappings.MergeMappings(&spec, &mapping)
		Expect(result).To(Not(BeNil()))
		Expect(result.NetworkMappings).To(BeNil())
		Expect(result.StorageMappings).To(BeNil())
	})
	table.DescribeTable("should merge the mappings ", func(primaryMapping *[]v2vv1alpha1.ResourceMappingItem, secondaryMapping *[]v2vv1alpha1.ResourceMappingItem, expected *[]v2vv1alpha1.ResourceMappingItem) {
		mapping := v2vv1alpha1.OvirtMappings{
			NetworkMappings: primaryMapping,
			StorageMappings: primaryMapping,
		}

		externalMapping := v2vv1alpha1.OvirtMappings{
			NetworkMappings: secondaryMapping,
			StorageMappings: secondaryMapping,
		}

		spec := v2vv1alpha1.ResourceMappingSpec{
			OvirtMappings: &externalMapping,
		}

		result := mappings.MergeMappings(&spec, &mapping)
		Expect(*result.NetworkMappings).To(ConsistOf(*expected))
		Expect(*result.StorageMappings).To(ConsistOf(*expected))
	},
		table.Entry("Primary nil",
			nil,
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, nil)},
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, nil)}),
		table.Entry("Secondary nil",
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, nil)},
			nil,
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, nil)}),
		table.Entry("Both input slices empty",
			&[]v2vv1alpha1.ResourceMappingItem{},
			&[]v2vv1alpha1.ResourceMappingItem{},
			&[]v2vv1alpha1.ResourceMappingItem{}),
		table.Entry("Primary empty",
			&[]v2vv1alpha1.ResourceMappingItem{},
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, nil)},
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, nil)}),
		table.Entry("Secondary empty",
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, nil)},
			&[]v2vv1alpha1.ResourceMappingItem{},
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, nil)}),

		table.Entry("Primary item with all nil values empty",
			&[]v2vv1alpha1.ResourceMappingItem{i(nil, nil, nil)},
			&[]v2vv1alpha1.ResourceMappingItem{},
			&[]v2vv1alpha1.ResourceMappingItem{}),
		table.Entry("Secondary item with all nil values empty",
			&[]v2vv1alpha1.ResourceMappingItem{},
			&[]v2vv1alpha1.ResourceMappingItem{i(nil, nil, nil)},
			&[]v2vv1alpha1.ResourceMappingItem{}),

		table.Entry("Primary item with all nil values plus other, named item",
			&[]v2vv1alpha1.ResourceMappingItem{i(nil, nil, nil), i(&id1, &name1, nil)},
			&[]v2vv1alpha1.ResourceMappingItem{},
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, nil)}),

		table.Entry("Disjuntive mappings with id and name",
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, nil)},
			&[]v2vv1alpha1.ResourceMappingItem{i(&id2, &name2, nil)},
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, nil), i(&id2, &name2, nil)}),
		table.Entry("Disjuntive mappings with id ",
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, nil, nil)},
			&[]v2vv1alpha1.ResourceMappingItem{i(&id2, nil, nil)},
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, nil, nil), i(&id2, nil, nil)}),
		table.Entry("Disjuntive mappings with name",
			&[]v2vv1alpha1.ResourceMappingItem{i(nil, &name1, nil)},
			&[]v2vv1alpha1.ResourceMappingItem{i(nil, &name2, nil)},
			&[]v2vv1alpha1.ResourceMappingItem{i(nil, &name1, nil), i(nil, &name2, nil)}),
		table.Entry("Disjuntive mappings with primary: id-only and secondary: name-only",
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, nil, nil)},
			&[]v2vv1alpha1.ResourceMappingItem{i(nil, &name2, nil)},
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, nil, nil), i(nil, &name2, nil)}),
		table.Entry("Disjuntive mappings with primary: name-only and secondary: id-only",
			&[]v2vv1alpha1.ResourceMappingItem{i(nil, &name1, nil)},
			&[]v2vv1alpha1.ResourceMappingItem{i(&id2, nil, nil)},
			&[]v2vv1alpha1.ResourceMappingItem{i(nil, &name1, nil), i(&id2, nil, nil)}),

		table.Entry("Completely overlapping mappings with id and name",
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, &type1)},
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, &type2)},
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, &type1)}),
		table.Entry("Completely overlapping mappings with id",
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, nil, &type1)},
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, nil, &type2)},
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, nil, &type1)}),
		table.Entry("Completely overlapping mappings with name",
			&[]v2vv1alpha1.ResourceMappingItem{i(nil, &name1, &type1)},
			&[]v2vv1alpha1.ResourceMappingItem{i(nil, &name1, &type2)},
			&[]v2vv1alpha1.ResourceMappingItem{i(nil, &name1, &type1)}),

		table.Entry("Mapping overlapping only with name",
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, &type1)},
			&[]v2vv1alpha1.ResourceMappingItem{i(nil, &name1, &type2)},
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, &type1)}),

		table.Entry("More primary mappings with name and id",
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, &type1), i(&id2, &name2, &type1)},
			&[]v2vv1alpha1.ResourceMappingItem{i(&id2, &name2, &type2)},
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, &type1), i(&id2, &name2, &type1)}),

		table.Entry("More secondary mappings with name and id",
			&[]v2vv1alpha1.ResourceMappingItem{i(&id2, &name2, &type1)},
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, &type2), i(&id2, &name2, &type2)},
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, &type2), i(&id2, &name2, &type1)}),

		table.Entry("Overlapping mappings with same id and different names plus other primary mapping",
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, &type1), i(&id3, &name3, nil)},
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name2, &type2)},
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, &type1), i(&id3, &name3, nil)}),
		table.Entry("Overlapping mappings with same id and different names plus other secondary mapping",
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, &type1)},
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name2, &type2), i(&id3, &name3, nil)},
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, &type1), i(&id3, &name3, nil)}),

		table.Entry("Overlapping mappings with same name and different ids plus other primary mapping",
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, &type1), i(&id3, &name3, nil)},
			&[]v2vv1alpha1.ResourceMappingItem{i(&id2, &name1, &type2)},
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, &type1), i(&id2, &name1, &type2), i(&id3, &name3, nil)}),
		table.Entry("Overlapping mappings with same name and different ids plus other secondary mapping",
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, &type1)},
			&[]v2vv1alpha1.ResourceMappingItem{i(&id2, &name1, &type2), i(&id3, &name3, nil)},
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, &type1), i(&id2, &name1, &type2), i(&id3, &name3, nil)}),

		table.Entry("All-in-one pathological mapping",
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, &type1), i(&id3, nil, &type1), i(nil, &name4, &type1), i(nil, nil, nil)},
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name2, &type2), i(&id2, &name1, &type2), i(&id3, &name3, &type2), i(&id4, nil, &type2), i(nil, nil, nil)},
			&[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, &type1), i(&id2, &name1, &type2), i(&id3, nil, &type1), i(nil, &name4, &type1), i(&id4, nil, &type2)}),
	)
	It("Should merge mapping with only import CR mapping", func() {
		mapping := v2vv1alpha1.OvirtMappings{
			NetworkMappings: &[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, &type1)},
			StorageMappings: &[]v2vv1alpha1.ResourceMappingItem{i(&id2, &name2, &type2)},
		}
		result := mappings.MergeMappings(nil, &mapping)

		Expect(result).To(Not(BeNil()))
		Expect(*result.NetworkMappings).To(ConsistOf(*mapping.NetworkMappings))
		Expect(*result.StorageMappings).To(ConsistOf(*mapping.StorageMappings))
	})
	It("Should merge mapping with only import CR mapping - case II", func() {
		mapping := v2vv1alpha1.OvirtMappings{
			NetworkMappings: &[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, &type1)},
			StorageMappings: &[]v2vv1alpha1.ResourceMappingItem{i(&id2, &name2, &type2)},
		}
		spec := v2vv1alpha1.ResourceMappingSpec{
			OvirtMappings: nil,
		}

		result := mappings.MergeMappings(&spec, &mapping)

		Expect(result).To(Not(BeNil()))
		Expect(*result.NetworkMappings).To(ConsistOf(*mapping.NetworkMappings))
		Expect(*result.StorageMappings).To(ConsistOf(*mapping.StorageMappings))
	})
	It("Should merge mapping with only external CR mapping", func() {
		mapping := v2vv1alpha1.OvirtMappings{
			NetworkMappings: &[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, &type1)},
			StorageMappings: &[]v2vv1alpha1.ResourceMappingItem{i(&id2, &name2, &type2)},
		}
		spec := v2vv1alpha1.ResourceMappingSpec{
			OvirtMappings: &mapping,
		}
		result := mappings.MergeMappings(&spec, nil)

		Expect(result).To(Not(BeNil()))
		Expect(*result.NetworkMappings).To(ConsistOf(*mapping.NetworkMappings))
		Expect(*result.StorageMappings).To(ConsistOf(*mapping.StorageMappings))
	})
	It("Should merge network and storage mappings when both present", func() {
		mapping := v2vv1alpha1.OvirtMappings{
			NetworkMappings: &[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, &type1)},
			StorageMappings: &[]v2vv1alpha1.ResourceMappingItem{i(&id2, &name2, &type1)},
		}
		externalMapping := v2vv1alpha1.OvirtMappings{
			NetworkMappings: &[]v2vv1alpha1.ResourceMappingItem{i(&id3, &name3, &type2)},
			StorageMappings: &[]v2vv1alpha1.ResourceMappingItem{i(&id4, &name4, &type2)},
		}

		spec := v2vv1alpha1.ResourceMappingSpec{
			OvirtMappings: &externalMapping,
		}
		result := mappings.MergeMappings(&spec, &mapping)

		Expect(result).To(Not(BeNil()))
		Expect(*result.NetworkMappings).To(ConsistOf(i(&id1, &name1, &type1), i(&id3, &name3, &type2)))
		Expect(*result.StorageMappings).To(ConsistOf(i(&id2, &name2, &type1), i(&id4, &name4, &type2)))
	})
	It("Should merge network from import CR and storage from external CR", func() {
		mapping := v2vv1alpha1.OvirtMappings{
			NetworkMappings: &[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, &type1)},
		}
		externalMapping := v2vv1alpha1.OvirtMappings{
			StorageMappings: &[]v2vv1alpha1.ResourceMappingItem{i(&id4, &name4, &type2)},
		}

		spec := v2vv1alpha1.ResourceMappingSpec{
			OvirtMappings: &externalMapping,
		}
		result := mappings.MergeMappings(&spec, &mapping)

		Expect(result).To(Not(BeNil()))
		Expect(*result.NetworkMappings).To(ConsistOf(i(&id1, &name1, &type1)))
		Expect(*result.StorageMappings).To(ConsistOf(i(&id4, &name4, &type2)))
	})
	It("Should merge network from external CR and storage from import CR", func() {
		mapping := v2vv1alpha1.OvirtMappings{
			StorageMappings: &[]v2vv1alpha1.ResourceMappingItem{i(&id2, &name2, &type1)},
		}
		externalMapping := v2vv1alpha1.OvirtMappings{
			NetworkMappings: &[]v2vv1alpha1.ResourceMappingItem{i(&id3, &name3, &type2)},
		}

		spec := v2vv1alpha1.ResourceMappingSpec{
			OvirtMappings: &externalMapping,
		}
		result := mappings.MergeMappings(&spec, &mapping)

		Expect(result).To(Not(BeNil()))
		Expect(*result.NetworkMappings).To(ConsistOf(i(&id3, &name3, &type2)))
		Expect(*result.StorageMappings).To(ConsistOf(i(&id2, &name2, &type1)))
	})
	It("Should override network from external CR with import CR", func() {
		mapping := v2vv1alpha1.OvirtMappings{
			NetworkMappings: &[]v2vv1alpha1.ResourceMappingItem{i(&id3, &name3, &type2)},
		}
		externalMapping := v2vv1alpha1.OvirtMappings{
			NetworkMappings: &[]v2vv1alpha1.ResourceMappingItem{i(&id3, &name3, &type1)},
		}

		spec := v2vv1alpha1.ResourceMappingSpec{
			OvirtMappings: &externalMapping,
		}
		result := mappings.MergeMappings(&spec, &mapping)

		Expect(result).To(Not(BeNil()))
		Expect(*result.NetworkMappings).To(ConsistOf(i(&id3, &name3, &type2)))
	})
	It("Should override storage from external CR with import CR", func() {
		mapping := v2vv1alpha1.OvirtMappings{
			StorageMappings: &[]v2vv1alpha1.ResourceMappingItem{i(&id4, &name4, &type2)},
		}
		externalMapping := v2vv1alpha1.OvirtMappings{
			StorageMappings: &[]v2vv1alpha1.ResourceMappingItem{i(&id4, &name4, &type1)},
		}

		spec := v2vv1alpha1.ResourceMappingSpec{
			OvirtMappings: &externalMapping,
		}
		result := mappings.MergeMappings(&spec, &mapping)

		Expect(result).To(Not(BeNil()))
		Expect(*result.StorageMappings).To(ConsistOf(i(&id4, &name4, &type2)))
	})
	It("Should merge and override network and storage mappings when both present", func() {
		mapping := v2vv1alpha1.OvirtMappings{
			NetworkMappings: &[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, &type1), i(&id3, &name3, &type1)},
			StorageMappings: &[]v2vv1alpha1.ResourceMappingItem{i(&id2, &name2, &type1)},
		}
		externalMapping := v2vv1alpha1.OvirtMappings{
			NetworkMappings: &[]v2vv1alpha1.ResourceMappingItem{i(&id1, &name1, &type2)},
			StorageMappings: &[]v2vv1alpha1.ResourceMappingItem{i(&id2, &name2, &type2), i(&id4, &name4, &type2)},
		}

		spec := v2vv1alpha1.ResourceMappingSpec{
			OvirtMappings: &externalMapping,
		}
		result := mappings.MergeMappings(&spec, &mapping)

		Expect(result).To(Not(BeNil()))
		Expect(*result.NetworkMappings).To(ConsistOf(i(&id1, &name1, &type1), i(&id3, &name3, &type1)))
		Expect(*result.StorageMappings).To(ConsistOf(i(&id2, &name2, &type1), i(&id4, &name4, &type2)))
	})
})

func i(id *string, name *string, tp *string) v2vv1alpha1.ResourceMappingItem {
	return v2vv1alpha1.ResourceMappingItem{
		Source: v2vv1alpha1.Source{
			ID:   id,
			Name: name,
		},
		Type: tp,
	}
}

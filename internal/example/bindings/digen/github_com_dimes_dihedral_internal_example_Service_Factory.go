// Code generated by go generate; DO NOT EDIT.
package digen

import target_pkg "github.com/dimes/dihedral/internal/example"

func factory_github_com_dimes_dihedral_internal_example_Service(d *DihedralServiceComponent) (*target_pkg.Service, error) {
	target := &target_pkg.Service{}
	ServiceTimeout, err := d.provides_github_com_dimes_dihedral_internal_example_ServiceTimeout()
	if err != nil {
		var zeroValue *target_pkg.Service
		return zeroValue, err
	}
	target.ServiceTimeout = ServiceTimeout
	DBStore, err := factory_github_com_dimes_dihedral_internal_example_dbstore_MemoryDBStore(d)
	if err != nil {
		var zeroValue *target_pkg.Service
		return zeroValue, err
	}
	target.DBStore = DBStore
	return target, nil
}


## Get started

```shell
kubectl apply -f config
```

Convert a CRD from YAML to JSON:
```shell
yq -j eval ../osf2f/config/crd/bases/osf2f.my.domain_profiles.yaml > crds/profiles.json
```
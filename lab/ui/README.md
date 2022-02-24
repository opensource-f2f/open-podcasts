
Demo address: http://103.61.38.146:30001/

## Get started

```shell
npm run start1
```

### In Kubernetes

```shell
kubectl apply -f config
```

Convert a CRD from YAML to JSON:
```shell
yq -j eval ../osf2f/config/crd/bases/osf2f.my.domain_episodes.yaml > crds/profiles.json
yq -j eval ../osf2f/config/crd/bases/osf2f.my.domain_profiles.yaml > crds/profiles.json
yq -j eval ../osf2f/config/crd/bases/osf2f.my.domain_rsses.yaml > crds/rsses.json
```

## Want to help?

This project uses [React](https://reactjs.org/) as the front-end framework, [Express](https://expressjs.com/) as the backend router.

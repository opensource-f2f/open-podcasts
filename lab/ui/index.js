const express = require('express');
const k8s = require("@kubernetes/client-node");
const {Client} = require("@kubernetes/client-node");
const crd = require("./crds/episodes.json"); //Import the express dependency
const app = express();              //Instantiate an express app, the main work horse of this server
const bodyParser = require('body-parser');
app.use(bodyParser());
const port = 5000;                  //Save the port number where your server will be listening

//Idiomatic expression in express to route and respond to a client request
app.get('/', (req, res) => {        //get requests to the root ("/") will route here
    res.sendFile('index.html', {root: __dirname});      //server responds by sending the index.html file to the client's browser
                                                        //the .sendFile method needs the absolute path to the file, see: https://expressjs.com/en/4x/api.html#res.sendFile

});

app.get('/a', (req, res) => {
    const k8s = require('@kubernetes/client-node');

    const kc = new k8s.KubeConfig();
    kc.loadFromDefault();

    const k8sApi = kc.makeApiClient(k8s.CoreV1Api);
    result = k8sApi.readNamespace("default")
    result.then(response => {
        res.end(response.body.metadata.name)
    })
});

app.get('/rsses', (req, res) => {
    const Client = require('kubernetes-client').Client
    const crd = require('./crds/rsses.json')
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crd)

    const all = client.apis['osf2f.my.domain'].v1alpha1.namespaces('default').rsses.get()
    all.then(response => {
        res.end(JSON.stringify(response.body.items))
    })
})

app.post('/rsses', (req, res) => {
    const Client = require('kubernetes-client').Client
    const crd = require('./crds/rsses.json')
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crd)

    client.apis['osf2f.my.domain'].v1alpha1.namespaces('default').rsses.post({
        body: {
            apiVersion: 'osf2f.my.domain/v1alpha1',
            kind: 'RSS',
            metadata: {
                name: req.body.name,
            },
            spec: {
                address: req.body.address
            }
        }
    })
    res.redirect('/')
})

app.get('/episodes', (req, res) => {
    const Client = require('kubernetes-client').Client
    const crd = require('./crds/episodes.json')
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crd)

    const all = client.apis['osf2f.my.domain'].v1alpha1.namespaces('default').episodes.get()
    all.then(response => {
        res.end(JSON.stringify(response.body.items))
    })
});

app.listen(port, () => {            //server starts listening for any attempts from a client to connect at port: {port}
    console.log(`Now listening on port ${port}`);
});

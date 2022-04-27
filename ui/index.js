const bodyParser = require('body-parser');
const { Client } = require('kubernetes-client');
const express = require('express');
const { http: https } = require('follow-redirects');
const { toXML } = require("to-xml");
const YAML = require('yaml');
const compare = require('compare')

const crdCategories = require('./crds/categories.json')
const crdEpisodes = require('./crds/episodes.json')
const crdNotifiers = require('./crds/notifiers.json')
const crdProfiles = require('./crds/profiles.json')
const crdRSSES = require('./crds/rsses.json') //Import the express dependency
const crdSubscriptions = require('./crds/subscriptions.json')

const port = 5000; //Save the port number where your server will be listening

const app = express(); //Instantiate an express app, the main work horse of this server
app.use(bodyParser());

app.use(express.static('build'))
//Idiomatic expression in express to route and respond to a client request
app.get('/', (req, res) => { //get requests to the root ("/") will route here
    res.sendFile('build/index.html', { root: __dirname }); //server responds by sending the index.html file to the client's browser
});

app.get('/healthz', (req, res) => { //get requests to the root ("/") will route here
    res.end('ok')
});

app.get('/stream/*', (req, res) => {
    try {
        const targetURL = req.url.replaceAll('/stream/', 'http://')
        https.get(targetURL, (rsp) => {
            for (let item in rsp.headers) {
                res.setHeader(item, rsp.headers[item])
            }
            rsp.pipe(res)
        })
    } catch (e) {
        console.error(e)
    }
})

const commandArgs = require('minimist')(process.argv.slice(2))
if (!commandArgs['defaultNamespace'] || commandArgs['defaultNamespace'] === "") {
    commandArgs['defaultNamespace'] = 'osf2f-system'
}
const defaultNamespace = commandArgs['defaultNamespace']

app.get('/namespaces', (req, res) => {
    const client = new Client({ version: '1.13' })

    client.api.v1.namespaces.get().then(resp => {
        res.end(JSON.stringify(resp.body.items, undefined, 2))
    })
})

app.get('/rsses', (req, res) => {
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crdRSSES)

    client
        .apis['osf2f.my.domain']
        .v1alpha1.namespaces(defaultNamespace)
        .rsses
        .get()
        .then(response => {
            var space = 0
            if (req.query.pretty === 'true') {
                space = 2
            }
            res.end(JSON.stringify(response.body.items, undefined, space))
        })
})

app.get('/rsses/search', (req, res) => {
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crdRSSES)
    const category = req.query.category

    client
        .apis['osf2f.my.domain']
        .v1alpha1
        .namespaces(defaultNamespace)
        .rsses
        .get()
        .then(response => {
            let items = response.body.items
            items = items.filter(function (val, index, arr) {
                if (val.spec.categories) {
                    const result = val.spec.categories.filter(function (categoryName) {
                        return category.toLowerCase() === categoryName.toLowerCase()
                    })
                    return result.length > 0
                }
                return false
            })

            var space = 0
            if (req.query.pretty === 'true') {
                space = 2
            }
            res.end(JSON.stringify(items, undefined, space))
        })
})

app.get('/rsses/subscribed', (req, res) => {
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crdRSSES)
    client.addCustomResourceDefinition(crdProfiles)
    client.addCustomResourceDefinition(crdSubscriptions)
    const profileName = req.query.profile
    var space = 0
    if (req.query.pretty === 'true') {
        space = 2
    }

    const rsses = []
    client
        .apis['osf2f.my.domain']
        .v1alpha1
        .namespaces(defaultNamespace)
        .profiles(profileName)
        .get()
        .then(response => {
            const profile = response.body
            if (profile.spec && profile.spec.subscription && profile.spec.subscription.name) {
                const subName = profile.spec.subscription.name
                client
                    .apis['osf2f.my.domain']
                    .v1alpha1
                    .namespaces(defaultNamespace)
                    .subscriptions(subName)
                    .get()
                    .then(async rsp => {
                        if (rsp.body.spec && rsp.body.spec.rssList) {
                            for (let i = 0; i < rsp.body.spec.rssList.length; i++) {
                                const item = rsp.body.spec.rssList[i]
                                const rss = await client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).rsses(item.name).get()
                                rsses.push(rss.body)
                            }
                            res.end(JSON.stringify(rsses, undefined, space))
                        }
                    })
            } else {
                res.end(JSON.stringify(rsses, undefined, space))
            }
        })
})

app.get('/rsses/one', (req, res) => {
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crdRSSES)
    const name = req.query.name

    client
        .apis['osf2f.my.domain']
        .v1alpha1
        .namespaces(defaultNamespace)
        .rsses(name)
        .get()
        .then(response => {
            var space = 0
            if (req.query.pretty === 'true') {
                space = 2
            }
            res.end(JSON.stringify(response.body, undefined, space))
        })
})

app.get('/rsses/export', (req, res) => {
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crdRSSES)

    client
        .apis['osf2f.my.domain']
        .v1alpha1
        .namespaces(defaultNamespace)
        .rsses
        .get()
        .then(response => {
            const items = response.body.items
            const exportItems = []
            for (var i = 0; i < items.length; i++) {
                exportItems[i] = {
                    apiVersion: items[i].apiVersion,
                    kind: "RSS",
                    metadata: {
                        name: items[i].metadata.name,
                    },
                    spec: {
                        address: items[i].spec.address,
                    }
                }
            }

            res.set({ "Content-Disposition": "attachment; filename=rsses.yaml" });

            var result = ""
            for (var i = 0; i < exportItems.length; i++) {
                const doc = new YAML.Document();
                doc.directivesEndMarker = true
                doc.contents = exportItems[i]
                result += doc.toString()
            }
            res.send(result);
        })
})

app.get('/rsses/opml/export', (req, res) => {
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crdRSSES)

    client
        .apis['osf2f.my.domain']
        .v1alpha1
        .namespaces(defaultNamespace)
        .rsses
        .get()
        .then(response => {
            const items = response.body.items

            const bodyItems = []
            for (var i = 0; i < items.length; i++) {
                bodyItems.push({
                    "@text": items[i].spec.title,
                    "@title": items[i].spec.title,
                    "@type": "rss",
                    "@xmlUrl": items[i].spec.address,
                    "@htmlUrl": ""
                })
            }
            const opml = {
                "?xml version='1.0' encoding='UTF-8' standalone='no' ?": null,
                "opml": {
                    "@version": "2.0",
                    head: {
                        title: 'From Open Podcasts',
                        dateCreated: new Date().toISOString(),
                    },
                    body: {
                        outline: bodyItems
                    }
                }
            }

            res.set({ "Content-Disposition": "attachment; filename=rsses.opml" });

            res.send(toXML(opml, null, 2));
        })
})

app.post('/rsses', (req, res) => {
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crdRSSES)

    client
        .apis['osf2f.my.domain']
        .v1alpha1
        .namespaces(defaultNamespace)
        .rsses
        .post({
            body: {
                apiVersion: 'osf2f.my.domain/v1alpha1',
                kind: 'RSS',
                metadata: {
                    generateName: (Math.random() + 1).toString(36).substring(7),
                },
                spec: {
                    address: req.body.address
                }
            }
        })
    res.end('ok')
})

app.get('/episodes', (req, res) => {
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crdEpisodes)
    client.addCustomResourceDefinition(crdRSSES)
    const rss = req.query.rss

    client
        .apis['osf2f.my.domain']
        .v1alpha1
        .namespaces(defaultNamespace)
        .rsses(rss)
        .get()
        .then(response => {
            const rssObj = response.body
            if (rssObj.status && rssObj.status.lastUpdateTime) {
                res.set({
                    'Last-Modified': rssObj.status.lastUpdateTime
                })
            }

            client
                .apis['osf2f.my.domain']
                .v1alpha1
                .namespaces(defaultNamespace)
                .episodes
                .get({ qs: { labelSelector: "rss=" + rss } })
                .then(response => {
                    var space = 0
                    if (req.query.pretty === 'true') {
                        space = 2
                    }
                    res.end(JSON.stringify(response.body.items.sort(compare.compareRevert), undefined, space))
                })
        })
});

app.get('/episodes/one', (req, res) => {
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crdEpisodes)
    const name = req.query.name

    client
        .apis['osf2f.my.domain']
        .v1alpha1
        .namespaces(defaultNamespace)
        .episodes(name)
        .get()
        .then(response => {
            var space = 0
            if (req.query.pretty === 'true') {
                space = 2
            }
            res.end(JSON.stringify(response.body, undefined, space))
        })
})

app.get('/profiles', (req, res) => {
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crdProfiles)
    const name = req.query.name

    var space = 0
    if (req.query.pretty === 'true') {
        space = 2
    }

    try {
        client
            .apis['osf2f.my.domain']
            .v1alpha1
            .namespaces(defaultNamespace)
            .profiles(name)
            .get()
            .then(response => {
                res.end(JSON.stringify(response.body, undefined, space))
            }).catch(err => {
                console.error('Error: ', err)
                res.send(JSON.stringify('{message:"error"}', undefined, space))
            })
    } catch (e) {
        console.error('Error: ', e)
        res.send(JSON.stringify('{message:"error"}', undefined, space))
    }
});

app.post('/profiles/create', (req, res) => {
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crdProfiles)
    const name = req.query.name

    try {
        client
            .apis['osf2f.my.domain']
            .v1alpha1
            .namespaces(defaultNamespace)
            .profiles
            .post({
                body: {
                    apiVersion: 'osf2f.my.domain/v1alpha1',
                    kind: 'Profile',
                    metadata: {
                        name: name,
                    }
                }
            }).error(res => {
                console.error(res)
            })
        res.end('ok')
    } catch (e) {
        console.error('Error: ', e)
    }
})

app.delete('/profile/playLater', (req, res) => {
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crdProfiles)

    client
        .apis['osf2f.my.domain']
        .v1alpha1
        .namespaces(defaultNamespace)
        .profiles(req.query.name)
        .get()
        .then(response => {
            var found = false
            targetProfile = response.body
            if (!targetProfile.spec) {
                targetProfile.spec = {}
            }
            if (!targetProfile.spec.laterPlayList) {
                targetProfile.spec.laterPlayList = []
            }

            targetProfile.spec.laterPlayList.forEach(function (item, index) {
                if (item.name === req.query.episode) {
                    targetProfile.spec.laterPlayList.splice(index, 1)
                    found = true
                    return false
                }
            })
            if (found) {
                client
                    .apis['osf2f.my.domain']
                    .v1alpha1
                    .namespaces(defaultNamespace)
                    .profiles(req.query.name)
                    .put({
                        body: targetProfile,
                    })
            }
        })
    res.status(200);
    res.end('ok')
})

app.post('/profile/playLater', (req, res) => {
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crdProfiles)

    client
        .apis['osf2f.my.domain']
        .v1alpha1
        .namespaces(defaultNamespace)
        .profiles(req.query.name)
        .get()
        .then(response => {
            var found = false
            targetProfile = response.body
            if (!targetProfile.spec) {
                targetProfile.spec = {}
            }
            if (!targetProfile.spec.laterPlayList) {
                targetProfile.spec.laterPlayList = []
            }

            targetProfile
                .spec
                .laterPlayList
                .forEach(function (item, index) {
                    if (item.name === req.query.episode) {
                        found = true
                        return false
                    }
                })
            if (!found) {
                targetProfile.spec.laterPlayList.push({
                    name: req.query.episode
                })

                client
                    .apis['osf2f.my.domain']
                    .v1alpha1
                    .namespaces(defaultNamespace)
                    .profiles(req.query.name)
                    .put({
                        body: targetProfile,
                    })
            }
        })
    res.status(200);
    res.end('ok')
})

app.post('/profile/playOver', (req, res) => {
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crdProfiles)

    const name = req.query.name
    const episode = req.query.episode

    client
        .apis['osf2f.my.domain']
        .v1alpha1
        .namespaces(defaultNamespace)
        .profiles(name)
        .get()
        .then(response => {
            targetProfile = response.body
            targetProfile.spec.laterPlayList.forEach(function (item, index) {
                if (item.name === episode) {
                    targetProfile.spec.laterPlayList.splice(index, 1)
                    if (targetProfile.spec.watchedList) {
                        targetProfile.spec.watchedList.push({ name: item.name })
                    } else {
                        targetProfile.spec.watchedList = [{
                            name: item.name
                        }]
                    }
                    return false
                }
            })

            client
                .apis['osf2f.my.domain']
                .v1alpha1
                .namespaces(defaultNamespace)
                .profiles(name)
                .put({
                    body: targetProfile,
                })
        })
    res.status(200);
    res.end('ok')
})

app.post('/profile/notifier', (req, res) => {
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crdProfiles)
    client.addCustomResourceDefinition(crdNotifiers)

    const name = req.query.name
    const url = req.query.url

    client
        .apis['osf2f.my.domain']
        .v1alpha1
        .namespaces(defaultNamespace)
        .profiles(name)
        .get()
        .then(response => {
            targetProfile = response.body
            if (!targetProfile.spec.notifier) {
                client
                    .apis['osf2f.my.domain']
                    .v1alpha1
                    .namespaces(defaultNamespace)
                    .notifiers
                    .post({
                        body: {
                            apiVersion: 'osf2f.my.domain/v1alpha1',
                            kind: 'Notifier',
                            metadata: {
                                generateName: name,
                            },
                            spec: {
                                feishu: {
                                    webhook_url: url,
                                },
                            }
                        }
                    }).then(res => {
                        targetProfile.spec.notifier = {
                            name: res.body.metadata.name,
                        }
                        client
                            .apis['osf2f.my.domain']
                            .v1alpha1
                            .namespaces(defaultNamespace)
                            .profiles(name)
                            .put({
                                body: targetProfile,
                            })
                    })
            } else {
                const notifierName = targetProfile.spec.notifier.name
                client
                    .apis['osf2f.my.domain']
                    .v1alpha1
                    .namespaces(defaultNamespace)
                    .notifiers(notifierName)
                    .get()
                    .then(notifierRes => {
                        const notifierObj = notifierRes.body
                        notifierObj.spec.feishu.webhook_url = url

                        client
                            .apis['osf2f.my.domain']
                            .v1alpha1
                            .namespaces(defaultNamespace)
                            .notifiers(notifierName)
                            .put({
                                body: notifierObj,
                            })
                    })
            }
        })
    res.status(200);
    res.end('ok')
})

app.get('/notifiers/one', (req, res) => {
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crdNotifiers)
    const name = req.query.name

    client
        .apis['osf2f.my.domain']
        .v1alpha1
        .namespaces(defaultNamespace)
        .notifiers(name)
        .get()
        .then(response => {
            var space = 0
            if (req.query.pretty === 'true') {
                space = 2
            }
            res.end(JSON.stringify(response.body, undefined, space))
        }).catch(err => {
            console.error('Error: ', err)
            res.send(JSON.stringify('{message:"error"}', undefined, space))
        })
})

app.post('/profile/social', (req, res) => {
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crdProfiles)

    const name = req.query.name
    const kind = req.query.kind
    const account = req.query.account

    client
        .apis['osf2f.my.domain']
        .v1alpha1
        .namespaces(defaultNamespace)
        .profiles(name)
        .get()
        .then(response => {
            targetProfile = response.body
            if (!targetProfile.spec.socialLinks) {
                targetProfile.spec.socialLinks = {}
            }
            targetProfile.spec.socialLinks[kind] = account

            client
                .apis['osf2f.my.domain']
                .v1alpha1
                .namespaces(defaultNamespace)
                .profiles(name)
                .put({
                    body: targetProfile,
                })
        })
    res.status(200);
    res.end('ok')
})

app.post('/subscribe', (req, res) => {
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crdProfiles)
    client.addCustomResourceDefinition(crdSubscriptions)

    const profileName = req.query.profile
    const rss = req.query.rss

    client
        .apis['osf2f.my.domain']
        .v1alpha1
        .namespaces(defaultNamespace)
        .profiles(profileName)
        .get()
        .then(response => {
            targetProfile = response.body
            if (!targetProfile.spec.subscription || !targetProfile.spec.subscription.name) {
                client
                    .apis['osf2f.my.domain']
                    .v1alpha1
                    .namespaces(defaultNamespace)
                    .subscriptions
                    .post({
                        body: {
                            apiVersion: 'osf2f.my.domain/v1alpha1',
                            kind: 'Subscription',
                            metadata: {
                                generateName: profileName,
                            },
                            spec: {
                                rssList: [{
                                    name: rss,
                                }],
                            }
                        }
                    }).then(res => {
                        targetProfile.spec.subscription = {
                            name: res.body.metadata.name,
                        }
                        client
                            .apis['osf2f.my.domain']
                            .v1alpha1
                            .namespaces(defaultNamespace)
                            .profiles(profileName)
                            .put({
                                body: targetProfile,
                            })
                    })
            } else {
                const subName = targetProfile.spec.subscription.name
                client
                    .apis['osf2f.my.domain']
                    .v1alpha1
                    .namespaces(defaultNamespace)
                    .subscriptions(subName)
                    .get()
                    .then(subRes => {
                        subscription = subRes.body
                        console.log(subscription + "-----" + subName)
                        if (subscription.spec.rssList) {
                            subscription.spec.rssList.push({
                                name: rss,
                            })
                        } else {
                            subscription.spec.rssList = [{
                                name: rss,
                            }]
                        }

                        client
                            .apis['osf2f.my.domain']
                            .v1alpha1
                            .namespaces(defaultNamespace)
                            .subscriptions(subName)
                            .put({
                                body: subscription,
                            })
                    })
            }
        })
    res.status(200);
    res.end('ok')
})

app.post('/unsubscribe', (req, res) => {
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crdProfiles)
    client.addCustomResourceDefinition(crdSubscriptions)

    const profileName = req.query.profile
    const rss = req.query.rss

    client
        .apis['osf2f.my.domain']
        .v1alpha1
        .namespaces(defaultNamespace)
        .profiles(profileName)
        .get()
        .then(response => {
            targetProfile = response.body
            if (targetProfile.spec.subscription) {
                const subName = targetProfile.spec.subscription.name
                client
                    .apis['osf2f.my.domain']
                    .v1alpha1
                    .namespaces(defaultNamespace)
                    .subscriptions(subName)
                    .get()
                    .then(subRes => {
                        subscription = subRes.body
                        if (subscription.spec.rssList) {
                            let found = false
                            subscription.spec.rssList.forEach(function (item, index) {
                                if (item.name === rss) {
                                    subscription.spec.rssList.splice(index, 1)
                                    found = true
                                    return false
                                }
                            })

                            if (found) {
                                client
                                    .apis['osf2f.my.domain']
                                    .v1alpha1
                                    .namespaces(defaultNamespace)
                                    .subscriptions(subName)
                                    .put({
                                        body: subscription,
                                    })
                            }
                        }
                    })
            }
        })
    res.status(200);
    res.end('ok')
})

app.get('/subscriptions/one', (req, res) => {
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crdSubscriptions)
    const name = req.query.name

    client
        .apis['osf2f.my.domain']
        .v1alpha1
        .namespaces(defaultNamespace)
        .subscriptions(name)
        .get()
        .then(response => {
            var space = 0
            if (req.query.pretty === 'true') {
                space = 2
            }
            res.end(JSON.stringify(response.body, undefined, space))
        })
})

app.get('/categories', (req, res) => {
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crdCategories)
    const name = req.query.name

    client
        .apis['osf2f.my.domain']
        .v1alpha1
        .categories()
        .get()
        .then(response => {
            var space = 0
            if (req.query.pretty === 'true') {
                space = 2
            }
            res.end(JSON.stringify(response.body, undefined, space))
        })
})

app.listen(port, () => { //server starts listening for any attempts from a client to connect at port: {port}
    console.log(`Now listening on port ${port}`);
});

# Wombat Builder

Wombat Builder is a tool for building distributions based on RedPanda Benthos. While from the outside it looks like a
simple website with some links, underneath there is a lot more going on to support the future ideas we have.

While it is easy to interpret this effort as a response to recent acquisition of Benthos by RedPanda, it is actually a
continuation of the ideas that were started much earlier; to create an ecosystem of components building on top of the
RedPanda Benthos core.

We decided the best way to support this would be through a build tool that allows you to cherry-pick the components you
want to use, and then build a distribution that includes them. This is the first step in a larger plan to create allow
more people to contribute to the ecosystem and enjoy the power of RedPanda Benthos.

## Help!
So by this point you must be wondering; "Hey, This is a funny project. How can I help?". Well, I am glad you asked. 

Helping out is actually quite easy. Just donate large sums of money to the project and we will be very happy! No? Well,
can't blame a man for trying.

If you are interested in helping out, please reach out to us in the [Wombat channel in the Gophers Slack](https://join.slack.com/t/gophers/shared_invite/zt-2k4y0h1c7-2YwUsmfllCc1utRRVM4qrA). We are always looking for people
who are interested in helping out with the project. We are also looking for people who are interested in creating
components that can be included in the distribution.

## How it works
This project is actually built on top of 3 different components; *the api*, *the service* and *the builders*. Let's 
start with the latter, shall we?

### The Builders
The builders are responsible for building the actual artifacts. They are capable of building wombat binaries based on
a `build` as it is provided to them by the service. The builders are also responsible for uploading the artifacts to
the correct location.

Builders listen for changes to the build definitions stored in the `builds` Nats Jetstream KV. When a change happens, every builder 
will be notified and will make an effort to 'claim' the build. Only one however will be successful in doing so. This
builder will then start the build process and update the status of the build in the KV. Once the build is complete, the
artifact is stored in the `artifacts` object store in Nats Jetstream.

As hinted, many different builders can be running at the same time, each with a different amount of workers associated.
This allows us to scale the build process horizontally, and to build many different artifacts at the same time.

### The Service
The service is a rather simple Nats Micro service that exposes a Request/Response API to manage builds. It contains
an endpoint `build.request` to which a request can be sent to create a new build. Unless an artifact is already being 
built for the given build configuration (os, arch, go version and packages), the service will update the build 
definition which in turn will notify the builders. Isn't that neat?

The service also keeps an internal search index which allows you to search for builds based on the build configuration.
Another endpoint is exposed for this purpose; `build.list`.

All service endpoints contain metadata describing what they do and what the data they require looks like. This metadata
can be consulted using the `nats micro ...` commands.

Services are also scalable components, and many of them can be spawned to handle the load of incoming requests.

### The API
While having a Nats Micro service will certainly give you a warm and fuzzy feeling, most people are still used to work
with a REST API. Therefore, we have created a simple API that sits on top of the service and translates REST requests
to Nats Micro requests. This API is also capable of serving the static content of the website.

Yes, this api is under-document. Well, it is actually not documented at all. The main reason for that is that it is 
still in flux and we are not sure yet what the final API will look like. We are also not sure if we will keep it at all.

The UI itself deserves a special notion. It is written in Vue3 using Vuetify and embedded into the go binary. Yes, I am
particularly proud of that. No seriously, I am a data geek and it took me to get it to this point.

# Disclaimer
This project is not affiliated with RedPanda and in no way are we trying to impersonate them. We are just a bunch of
people who like the project and want to see it grow. We are also not trying to make any money from this project. In
fact, we are spending money on it. 
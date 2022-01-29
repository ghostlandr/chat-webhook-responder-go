# chat-webhook-responder-go
Go module for chat webhook responder app

## Deployment

With gcloud sdk and the app-engine-go libs installed, run the following from the root:

```
gcloud app deploy . -v optional-version-name -q
```

If you want to deploy the default module, omit the -v param. If you want it to confirm first,
omit the -q param.

Some day I will automate this, but until then, this README entry will do!

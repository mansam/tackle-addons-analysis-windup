Tackle Add-ons - Analysis - Windup
---

This add-on analyzes the source code repository to find cloud-readiness issues in the code using [Windup](https://github.com/windup/windup).

## Building

We need to download the MTA CLI from https://developers.redhat.com/products/mta/download.
We then need to unzip it in the current folder and rename it to mta-cli.

```
podman build -t quay.io/konveyor/tackle-addons-analysis-windup:latest .
```


```
podman push quay.io/konveyor/tackle-addons-analysis-windup:latest
```

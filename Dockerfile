FROM scratch

COPY humioctl /humioctl

ENTRYPOINT ["/humioctl"]

FROM alpine

COPY rdocker /bin/rdocker

ENTRYPOINT ["/bin/rdocker"]
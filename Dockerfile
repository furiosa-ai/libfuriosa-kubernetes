FROM ghcr.io/furiosa-ai/furiosa-smi:latest as smi

FROM golang:1.21.7-bookworm

# Copy hwloc binaries and libraries from the builder stage
COPY --from=smi /usr/local/lib/libfuriosa_smi.so /usr/local/lib/libfuriosa_smi.so
COPY --from=smi /usr/local/include/furiosa/furiosa_smi.h /usr/local/include/furiosa/furiosa_smi.h

# Configure env values
ENV C_INCLUDE_PATH /usr/local/include
ENV LD_LIBRARY_PATH usr/local/lib

WORKDIR $GOPATH

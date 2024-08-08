FROM registry.furiosa.ai/furiosa/furiosa-smi:latest as smi

FROM golang:1.21.7-bookworm

# Copy hwloc binaries and libraries from the builder stage
COPY --from=smi /usr/lib/x86_64-linux-gnu/libfuriosa_smi.so /usr/lib/x86_64-linux-gnu/libfuriosa_smi.so
COPY --from=smi /usr/include/furiosa/furiosa_smi.h /usr/include/furiosa/furiosa_smi.h
RUN ldconfig

WORKDIR $GOPATH

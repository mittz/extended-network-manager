FROM centos

RUN curl https://releases.rancher.com/install-docker/1.12.sh | sh

COPY ./extended-network-manager /extended-network-manager

CMD ["/extended-network-manager"]
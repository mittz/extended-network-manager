FROM scratch

COPY ./rancher-pipework /rancher-pipework

CMD ["rancher-pipework"]
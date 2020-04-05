FROM scratch
LABEL author="Alberto Bregliano"

EXPOSE 3000

COPY main /

ENTRYPOINT [ "/main" ]

CMD ["-redis","127.0.0.1:6379"]
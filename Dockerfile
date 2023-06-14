FROM golang:1.18.1

# WORKDIR /home/ratul/AESL/Test/Library_Vivek/pkg

# RUN cd /home/ratul/AESL/Test/Library_Vivek/pkg && go build -o ../library

WORKDIR /home
COPY ./pkg /home

RUN cd /home && go build -o library

CMD [ "/home/library" ]
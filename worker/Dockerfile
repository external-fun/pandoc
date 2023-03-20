FROM golang:1.18.1
RUN apt update && apt install -y \
    texlive-latex-base \
    texlive-fonts-recommended \
    texlive-fonts-extra \
    texlive-latex-extra \
    pandoc
ADD . /app/
WORKDIR /app
EXPOSE 8080
CMD ["go", "run", "./main.go"]
sudo apt-get update && sudo apt-get install -y --no-install-recommends g++ make autoconf automake libtool nasm wget
wget https://github.com/mozilla/mozjpeg/releases/download/v3.2-pre/mozjpeg-3.2-release-source.tar.gz &&
    tar -xvzf mozjpeg-3.2-release-source.tar.gz &&
    rm mozjpeg-3.2-release-source.tar.gz &&
    cd mozjpeg &&
    ./configure &&
    make install &&
    cd ../ && sudo rm -rf mozjpeg &&
    sudo ln -s /opt/mozjpeg/bin/jpegtran /usr/local/bin/jpegtran &&
    sudo ln -s /opt/mozjpeg/bin/cjpeg /usr/local/bin/cjpeg

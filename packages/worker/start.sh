echo "Installing firecracker"
firecracker_url="https://github.com/firecracker-microvm/firecracker/releases"
latest=$(basename $(curl -fsSLI -o /dev/null -w  %{url_effective} ${firecracker_url}/latest))
arch=`uname -m`
curl -L ${firecracker_url}/download/${latest}/firecracker-${latest}-${arch}.tgz | tar -xz
mv release-${latest}-$(uname -m)/firecracker-${latest}-$(uname -m) /usr/local/bin/firecracker
rm -rf release-${latest}-$(uname -m)

echo "Installing vm kernel"
kernel_url="https://s3.amazonaws.com/spec.ccfc.min/img/quickstart_guide/$arch/kernels/vmlinux.bin"
curl -fsSL -o ./linux.bin $kernel_url

echo "Installing cni plugins"
cni_url="https://github.com/containernetworking/plugins/releases/download/v1.1.1/cni-plugins-linux-amd64-v1.1.1.tgz"
sudo mkdir -p /opt/cni/bin
curl -L $cni_url | tar -xz -C /opt/cni/bin

git clone https://github.com/awslabs/tc-redirect-tap
cd tc-redirect-tap && make
sudo mv tc-redirect-tap /opt/cni/bin
cd .. && rm -rf tc-redirect-tap

sudo mkdir -p /etc/cni/conf.d
sudo cp fcnet.conflist /etc/cni/conf.d/

echo "Building Agent + Worker"
cd ../agent && ./build.sh
mv ./rootfs.ext4 ../worker/
cd ../worker

go build -o worker
sudo MAX_VMS=10 VCPUS=1 MEM_SIZE=128 DRIVE_PATH=./rootfs.ext4 KERNEL_PATH=./linux.bin ./worker

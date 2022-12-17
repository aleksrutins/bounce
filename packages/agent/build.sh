dd if=/dev/zero of=rootfs.ext4 bs=1M count=1000
mkfs.ext4 rootfs.ext4

IMG_ID=$(docker build -q .)
CONTAINER_ID=$(docker run -d $IMG_ID)

mkdir -p /tmp/rootfs
sudo mount rootfs.ext4 /tmp/rootfs
sudo docker cp $CONTAINER_ID:/ /tmp/rootfs
sudo umount /tmp/rootfs

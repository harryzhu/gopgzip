# gopgzip

## build
`./build.sh`

## zip a file PARALLELLY
```
./gopgzip zip --input=/home/harryzhu/test/abc.tar --output=/the/path/of/the/abc.tar.gz --thread=6 --level=9
```
### Performance:
221GB: zip to 100GB in 44 minutes

 - `--threads=6`: from 1 to maximum(your-total-cpu-core)
 - `--level=0|1|6|9`:
   - 0 - No Compression
   - 1 - Best Speed
   - 6 - Default Compression
   - 9 - Best Compression

## unzip a file
```
./gopgzip unzip --input=/home/harryzhu/test/abc.tar.gz --output=/the/path/of/the/abc.tar`
```

## tar a folder recursively
```
./gopgzip tar --input=/the/path/of/the/folder --output=/home/harryzhu/test/abc.tar
```
 you can add `--compression=0 | 1 | 2` to define:
 - 0 - No Compression
 - 1 - gzip Compression
 - 2 - zstd Compression

## untar a file
```
./gopgzip untar --input=/home/harryzhu/test/abc.tar --output=/the/path/of/the/folder
```

 you have to add `--compression=0|1|2` same as the [tar] above

## encrypt a file
```
./gopgzip encrypt --input=/home/harryzhu/test/abc.tar --output=/the/path/of/the/abc.tar.enc
```
 in your `/etc/profile`, add `export HARRYZHUENCRYPTKEY=Your-Password` to set your own PASSWORD
 or use `--password="Your-Password" --force` inline (NOT recommend)

## decrypt a file
```
./gopgzip decrypt --input=/home/harryzhu/test/abc.tar.enc --output=/the/path/of/the/abc.tar
```

## md5sum a file: for a small file
```
./gopgzip md5 --input=/home/harryzhu/test/abc.tar.gz [--output=/the/path/of/the/abc.tar.md5]
```

## b3sum a file: for a very large file
```
./gopgzip b3sum --input=/home/harryzhu/test/abc.tar.gz [--output=/the/path/of/the/abc.tar.b3sum]
```
# IMT Digest

IMT Digest is an utility to download arbitrary-sized files from an internet source, apply an _IMT hash_ to its contents and save the hash to a file.

It includes a throttle feature to limit the download speed, so it doesn't consume all the available bandwidth. 


## Installation

`$ go get github.com/FcoManueel/imtdigest` 


## Usage

The simplest way to use it is by providing a file location and a url in the following way:
 
`$ imtdigest -file "/tmp/imt-hash.txt" -url "https://www.google.com/"`

If you want the download to be throttled then you have to provide the `-rate` flag followed by an amount of bytes per second. 

Here you can find the full explanation of the available flags:

`$ imtdigest --help`
```
Usage of imtdigest:
  -file string
    	The path of the file where the output will be saved. (required)
  -rate int
    	Limits the max download rate. Units are in bytes/second. (optional)
  -url string
    	The URL from which the data will be fetched. (required)
```


## Additional comments

### On dependencies 
For the throttling feature I decided to go with a token bucket algorithm.  
While there's nothing in the standard library for that, the Go team has indeed developed [an implementation of it](https://godoc.org/golang.org/x/time/rate).  
While my aim was use as few external packages as possible, I considered fair to include that one given its nature and origin.

### On throttling
The throttling I applied works only at an application level. 
There are other factors that could affect the effective transference rate (e.g. data could still be transferred if there's space on the OS receive buffer for the given TCP connection, even if the application thread consuming the data is sleeping).
  
The throttling _might_ be good enough (for the purpose of the exercise and the allotted time anyways) since we are likely dealing with sufficiently large files, but there's definitely room for improvement.
    
The current implementation is specially bad for small file sizes and if I were to start again I would go with a custom/simpler implementation of the token bucket.

### On hosting
I added this project to my public GitHub for ease of installation and evaluation, but let me know if you prefer me to take it down or make it private.

### On hashing
This is fairly irrelevant for the intent of the exercise, but I found it interesting enough to share. The logic provided for the hashing function specifies that when `i=0` then `h[i-1]` should be considered to be 0.
 This however might be an issue for a hashing 
function, since it causes that every 0 byte flushes the hash (see [hash_test.go](./hash_test.go) for examples of this).

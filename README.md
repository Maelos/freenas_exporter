# FreeNAS Custom Exporter

6/8/18 - I got some really helpful feedback here https://groups.google.com/forum/#!topic/golang-nuts/0o7jMAXPBiE and am incorporating that now.  I also am reducing the test file greatly, but the past testing should be in the file history so it is not lost.

6/7/18 - As of now I have the first part of the command working, but I don't think the piping works within the exec.Command() creation.  Judging from the example on gobyexample.com for exec'ing processes, I may need to use the command's standard out, in, and error to manipulate the data appropriately.  It may also be wise to use something besides grep and awk which I am sure are available within Go's standard libraries.  At this point I feel like I have forgotten part of the class I took and or just need to invest in a book and relearn a bit, specifically on the reader and reader implementing packages.


I am trying to create a custom exporter for use in FreeNAS as the node_exporter does not send the correct CPU temperature (currently).  You may find my troubleshooting in these, and a few other, links:
https://groups.google.com/forum/#!topic/prometheus-users/MjA77maIz5o
https://github.com/prometheus/node_exporter/issues/945

This will be the base code.  Once I get this working I hope to add in temp metrics for the drives, then anything else I or the FreeNAS community wishes.  I have a working shell script for the temperatures and know the FreeNAS community has quite a few more scripts that can be translated over to Go or to otherwise work with Prometheus for a more sophisticated/fun monitoring system.  Any guidance, advice, etc you have is welcome.  All questions welcome too!

Status: Compiles, I believe with dependencies, but does not run as descriped in the google groups thread.  Not sure what is missing, but I will keep playing around with it.

Credit for much of the work goes to:
https://rsmitty.github.io/Prometheus-Exporters/, https://www.robustperception.io/setting-a-prometheus-counter/, and looking at the node_xporter of Prometheus, specifically cpu_freebsd

[16:39] <Nikon_NLG[m]> maelos: Don't you check if you have some specific flag for kernel, and you'll have correct temperature for your MB/CPU ?

Renamed this to custom exporter and fixed a few things
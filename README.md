# Chiller

Support service to back up the [deslauriers.world](https://deslauriers.world) site data.

*There are way easier ways to go about this, mostly this was an excercise in playing around with Go's unmarshalling and crud functions.* 

## Context:

Because the site is hosted on broken computers, and more importantly, because I am constantly fiddling with it, crashes and failures happen.  

*Learned some hard lessons I already knew about why external backups are important.*

## Function: Runs nigthly back ups

1. Logs into site. 
1. Uses read access to micro-services to pull db table data.
    1. It would be easier to call the tables one at a time, but I am doing giant nested json for two reasons:
        1. Learning exercise for navigating json w/ Go.
        1. Ultimately, I will move to asyncronus pub/sub, and thats how the data will be. 
1. Loads into back up database outside K8s cluster.
    * field level encryption will be applied by the called services.

## Day Two

1.  Build function to populate data back to active site in the event of cluster failure.  

## Day Three

1. Implement pub/sub pattern vs blocking.
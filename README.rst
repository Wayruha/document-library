Systems Engineering 2 - Assignment 1
====================================


Task description
----------------

Here I try to develop an online **document library**.
Users are presented with an input form, where they can submit *documents* (e.g., books, poems, recipes) along with *metadata* (e.g., author, mime type, ISBN).
For the sake of simplicity, they can view *all* stored documents on a single page.


Nginx
~~~~~

   1. ``nginx/Dockerfile`` completed and verified. It shows the landing page

HBase
~~~~~

   1. ``hbase`` is working as it was given


ZooKeeper
~~~~~~~~~

   1. `zookeeper`` is used to communicate between HBase and servers (grproxy and gserv) 

grproxy
~~~~~~~
   
   1. ``grproxy` can create node to hbase `zookeeper`` and set watch on it.
   2.  if ``gserve`` creates ephemeral child node in ZooKeeper's node then ``grproxy`` will get notifications

gserve
~~~~~~

   1. gserve can create ephemeral child node under grproxy's defined node
   2. Two instances *gserve1* and *gserve2* can run and create child nodes and write their own service_name:port as data which is used by grproxy to select gserve and get their addresses to communicate
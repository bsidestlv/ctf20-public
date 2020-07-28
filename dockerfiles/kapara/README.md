# Kapara Challenge.

# Challenge Description.

Kapara achi. 

/home/user is writable and readable so you could put your binaries there
and other won't read them.

Cheers.


# Flag
BSidesTLV2020{KaparaOnYouKaparaOnMeKaparaOnEveryoneInTheWorld}

# Build:
docker build . -t kapara
docker run -it -p 127.0.0.1:5555:5555 kapara:latest

# Login:
ssh user@ip -p 5555 
passwd: user 

/home/user is writable and executable only to allow people putting their code
and not let others read their binaries.

# Challenge Description:

/dev/kapara0 is exposed with world 777 permissions, it's a kapara device 
it contains a linked list and an object ID management system.

The code itself is quite simple, kapara0 exposes an IOCTL interface to 
create, add, modify and delete a Kapara object.

For example, if one wishes to create a Kapara object they should perform the
following operations:

```
#define KAPARA_BASE_IOCTL 1336
#define KAPARA_ADD_OBJ KAPARA_BASE_IOCTL+1
#define KAPARA_DEL_LIST_OBJECT KAPARA_BASE_IOCTL+2 
#define KAPARA_FIND_OBJ KAPARA_BASE_IOCTL+3
#define KAPARA_CREATE_OBJ KAPARA_BASE_IOCTL+4
#define KAPARA_INIT_DA KAPARA_BASE_IOCTL+5
#define KAPARA_DEL_OBJ KAPARA_BASE_IOCTL+6
#define KAPARA_FIND_OBJ_ID KAPARA_BASE_IOCTL+8
#define KAPARA_BKDR KAPARA_BASE_IOCTL+9



      typedef struct kapara_obj {
              int channel_id;
              char name[32];
              int popularity;
      }kapara_obj;

      #define ic(txt, ioctl_nm, pm, b) do { \
        puts(txt);\
        if (b) \
          getchar(); \
        ioctl(fd, ioctl_nm, pm); \
      } while(0);
      kapara_obj kob = {
        .channel_id = 0,
        .name = "hello",
        .popularity = 4321,
      };
      int id = 1337;
      int fd = open("/dev/kapara0", O_RDONLY);

      ic("setting up", KAPARA_INIT_DA , &kob, 0 )
      ic("getting id ", 1340, &kob , 1)
      ic("adding ", 1337, &kob.channel_id ,0 )
      for (int i = 0; i < 1; i++) {
        ic("pumping the slab", 1345, &b , 0);
      }
```
The kapara linked list operations are mutex protected, however the ID management 
operations, (peformed by Linux idr\* functions), are not protected by any kind of
locking mechanism. This creates an interesting scenario where one can delete
a kapara object from the idr list and free it, but the kapara linked list object
will still reference to it, allowing the attacker to exploit this behavior
and gain arbitrary read and write primitives.

I believe the hard parts in this challenge are as following:

1. Identifying the idr\* and kapara objects management interfaces, (reversing the ko).
2. Finding an Object that will fit to a kapara object size allocated in 
the same kernel heap slab.
3. Triggering the race reliably without crashing the system, i.e allocating the
object and filling it properly.
4. Creating an arbitrary r/w interface in order to elevate privileges.


Cheers,
  G

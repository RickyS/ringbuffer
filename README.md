Simple package that implements a ring buffer — a circular buffer, also known as a FIFO
=======================================================


To install:  
       $ go get github.com/RickyS/ringbuffer  

A ring buffer is an array where stuff is put into one end and removed from the other.  Also known as a queue.
The data is represented as an array of interface elements.

To Use:  
         type myThing ringbuffer.RingElement  // Create your element type as you need.  

         var whatever == myThing("whatever") // Assuming a conversion from string.  Not needed.  

         rb := RingBuffer.New(40)           // Create the fixed-size ringbuffer.  
         rb.Write(myThing) // Et cetera     // Insert the first element.  
         aThing := rb.Read()                // Retrieve the next element.  

         for rb.HasAny() {        // Classic loop to read out all the elements of the ringbuffer.  
             doSomethingWith(rb.Read())  
         }
         // Reading an empty ringbuffer returns nil.             

 THIS IS NOT CONCURRENT   
If you want a concurrent ringbuffer, consider using a channel.

The test suite is larger and more complex than the code being tested.

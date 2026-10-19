package route

import "sync"

// backlink_mutex serializes a document save together with its backlink index
// update. Holding it across both the content commit and the re-index removes the
// window in which a concurrent commit could slip in between saving content and
// writing backlinks. Only one holder at a time (single process); the holder owns
// the whole critical section so no other code re-enters it (no nested locking).
var backlink_mutex sync.Mutex

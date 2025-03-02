package chanUtil

func listener(listenCh, done chan struct{}, callback func()) {
	for true {
		select {
		case <-done:
			return
		case _, open := <-listenCh:
			if !open {
				return
			}
		}

		callback()
	}
}

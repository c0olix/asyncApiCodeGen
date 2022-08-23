package de.hochbahn.events;

import javax.validation.constraints.NotNull;

public interface IUserLockedEventProducer {
    void produceUserLockedEvent(@NotNull UserLockedEvent userLockedEvent);
}
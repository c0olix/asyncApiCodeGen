package de.hochbahn.events;

import java.util.function.Consumer;
import org.springframework.messaging.Message;

public interface IUserDeletedEventConsumer {
	Consumer<Message<UserDeletedEvent>> consumeUserDeletedEvent();
}
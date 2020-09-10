import React, { useState, useEffect } from "react";
import { Box, Flex, Text, useColorMode, Center } from "@chakra-ui/core";
import { useEventsStreamSubscription, Event, useSharedBoardQuery } from "@app/generated/graphql";
import EventCard from "@app/components/atoms/eventCard";
import Header from "../atoms/header";

interface ObserverBasicLayoutProps {
    initialTitle: string;
    initialMessage: string;
    bufferSize?: number;
}

const ObserverBasicLayout: React.FC<ObserverBasicLayoutProps> = ({
    initialTitle,
    initialMessage,
    bufferSize = 3 * 5,
}: ObserverBasicLayoutProps) => {
    const [events, setEvents] = useState<Event[]>([]);
    const { loading, error, data: sharedEvents } = useSharedBoardQuery({ variables: { size: bufferSize } });
    const sub = useEventsStreamSubscription();

    useEffect(() => {
        console.log(sub.data);
        if (sub.data?.eventStream) {
            const event = sub.data.eventStream as Event;
            if (events.length + 1 > bufferSize) {
                events.pop();
            }
            setEvents([event, ...events]);
        }
    }, [sub.data]);

    useEffect(() => {
        if (sharedEvents) {
            setEvents(sharedEvents.sharedBoard as Event[]);
        }
    }, [sharedEvents]);

    if (loading) {
        return (
            <Center>
                <Text>Loading...</Text>
            </Center>
        );
    }

    return (
        <Box>
            <Header initialTitle={initialTitle} initialMessage={initialMessage} />
            <Box
                mt={16}
                p={4}
                display={"grid"}
                gridAutoFlow={"row dense"}
                gridTemplateColumns={"repeat(auto-fit, minmax(250px, 350px))"}
                gridGap={"18px"}
                justifyItems={"center"}
                justifyContent={"center"}
            >
                {events.map((event: Event) => (
                    <EventCard key={event.id} event={event} />
                ))}
            </Box>
        </Box>
    );
};

export { ObserverBasicLayout };

import React, { useState, useEffect, ReactElement } from "react";
import { Box, Flex, Text, useColorMode, Center, Image } from "@chakra-ui/core";
import { useEventsStreamSubscription, Event, useSharedBoardQuery, useHeaderQuery } from "@app/generated/graphql";
import EventCard from "@app/components/atoms/eventCard";
import Header from "../atoms/header";
import SSLogo from "../atoms/logo";

interface ObserverBasicLayoutProps {
    initialTitle?: string;
    hashtag?: string;
    bufferSize?: number;
    brandData?: string;
}

const ObserverBasicLayout: React.FC<ObserverBasicLayoutProps> = ({ bufferSize = 25 }: ObserverBasicLayoutProps) => {
    const [events, setEvents] = useState<Event[]>([]);
    const { loading, error, data: sharedEvents } = useSharedBoardQuery({ variables: { size: bufferSize } });
    const sub = useEventsStreamSubscription();

    const { data: header } = useHeaderQuery();

    useEffect(() => {
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

    let currentBrand: string | ReactElement = header?.header.brand;

    if (header?.header.brand.startsWith("http")) {
        // if is link, it'll be used as a image
        currentBrand = <Image src={header?.header.brand} maxHeight={6} />;
    }

    if (loading) {
        return (
            <Center height={"100vh"}>
                <SSLogo width={80} height={80} />
            </Center>
        );
    }

    return (
        <Box>
            <Header brand={currentBrand} initialTitle={header?.header.title} hashtag={header?.header.hashtag} />
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

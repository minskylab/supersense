import React, { useState, useEffect } from "react";
// import { Box, Text, Flex, Image } from "theme-ui";
import type { Event } from "@app/generated/graphql";
import { GitHub, Twitter, Info, Share } from "react-feather";
import { format } from "timeago.js";
import { useSpring, animated } from "react-spring";
import { Flex, Text, Image, Box, useTheme, Center, useColorMode, useColorModeValue, Avatar } from "@chakra-ui/core";

interface EventCardProps {
    event: Event;
}

const EventCard: React.FC<EventCardProps> = ({ event }: EventCardProps) => {
    const [timeAgo, setTimeAgo] = useState<string>(format(event.emittedAt, "en_US"));
    const props = useSpring({ opacity: 1, from: { opacity: 0 } });
    const theme = useTheme();

    useEffect(() => {
        const timer = setInterval(() => {
            setTimeAgo(format(event.emittedAt, "en_US"));
        }, 1000);

        return () => {
            window.clearInterval(timer);
        };
    }, []);

    const eventKindToTitle = (eventKind: string): string => {
        eventKind[0].toUpperCase();
        const sPart = eventKind.substring(1).replace(/-/g, " ");
        return eventKind.substring(0, 1).toUpperCase() + sPart;
    };

    const borderColor = useColorModeValue(theme.colors.gray["900"], theme.colors.gray["600"]);
    const highLightColor = useColorModeValue(theme.colors.gray["50"], theme.colors.gray["700"]);

    return (
        <animated.div
            key={event.id}
            style={{
                ...props,
                display: "flex",
                flexDirection: "column",
                // padding: "1rem",
                width: "100%",
                border: "solid 2px " + borderColor,
                borderRadius: "8px",
                WebkitBoxShadow: "5px 5px 0px 0px " + borderColor,
                MozBoxShadow: "5px 5px 0px 0px " + borderColor,
                boxShadow: "5px 5px 0px 0px " + borderColor,
            }}
        >
            <Flex sx={{ alignContent: "center" }} py={3} px={3}>
                <Box>
                    {event.sourceName === "twitter" && (
                        <Twitter size={24} fill={theme.colors.cyan["400"]} strokeWidth={0} />
                    )}
                    {event.sourceName === "github" && (
                        <GitHub size={24} fill={theme.colors.gray["400"]} strokeWidth={0} />
                    )}
                    {event.sourceName === "dummy" && (
                        <Info size={24} fill={theme.colors.yellow["400"]} strokeWidth={0} />
                    )}
                    {event.sourceName === "spokesman" && (
                        <Info size={24} fill={theme.colors.pink["400"]} strokeWidth={0} />
                    )}
                </Box>
                <Box ml={2}>{eventKindToTitle(event.eventKind)}</Box>
                <Box ml={"auto"}>
                    <Text opacity={0.5} className="timeago">
                        {timeAgo}
                    </Text>
                </Box>
            </Flex>
            {event.entities.media.length > 0 ? <Image src={event.entities.media[0].url}></Image> : null}
            <Center py={4} px={4} my={"auto"}>
                <Box maxWidth={"100%"}>
                    <Text sx={{ textAlign: "center" }} as={"h3"}>
                        {event.message.slice(0, 220)}
                    </Text>
                </Box>
            </Center>
            <Flex sx={{ alignContent: "center" }} mt={"auto"} bg={highLightColor} borderBottomRadius={8} py={2} px={3}>
                <Avatar src={event.actor.photo} size={"sm"} />
                <Box
                    ml={3}
                    sx={{
                        display: "flex",
                        flexDirection: "column",
                        justifyContent: "center",
                    }}
                >
                    <Text as={"p"} sx={{ alignContent: "center" }}>
                        {event.actor.username || event.actor.name}
                    </Text>
                </Box>
                <Box
                    ml={"auto"}
                    sx={{
                        display: "flex",
                        flexDirection: "column",
                        justifyContent: "center",
                    }}
                >
                    <Box _hover={{ cursor: "pointer" }}>
                        <Share
                            onClick={() => {
                                window.open(event.shareURL, "_blank");
                            }}
                        />
                    </Box>
                </Box>
            </Flex>
        </animated.div>
    );
};

export default EventCard;

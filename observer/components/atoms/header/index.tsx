import React, { FC, ReactElement, Component } from "react";
import { Flex, Box, Text, useColorMode, useColorModeValue } from "@chakra-ui/core";
import { Settings, Moon, Sun } from "react-feather";

interface HeaderProps {
    brand?: ReactElement | string;
    initialTitle: string;
    hashtag: string;
    onSettings?: () => void;
    lightColor?: string;
    darkColor?: string;
}

const Header: FC<HeaderProps> = ({
    initialTitle,
    hashtag,
    onSettings,
    brand = "SUPERSENSE",
    lightColor = "teal.200",
    darkColor = "teal.700",
}: HeaderProps) => {
    const { colorMode, toggleColorMode } = useColorMode();

    const bgColor = useColorModeValue(lightColor, darkColor);

    const changeColorModeIcon = useColorModeValue(
        <Moon onClick={toggleColorMode} />,
        <Sun onClick={toggleColorMode} />
    );

    return (
        <Flex
            p={3}
            bg={bgColor}
            position={"fixed"}
            top={0}
            width={"100%"}
            zIndex={10}
            // borderBottomStyle={"solid"}
            // borderBottom={"1px"}
            // borderBottomColor={"black"}
        >
            <Box flex={1}>{brand}</Box>
            <Flex flex={[1, 2, 4]}>
                <Text display={["none", "block"]}>{initialTitle}</Text>
            </Flex>
            <Flex flex={1} justifyContent={"flex-end"}>
                <Text fontFamily={"body"}>{hashtag}</Text>
            </Flex>
            {/* <Flex ml={4}>
                <Settings onClick={onSettings} />
            </Flex> */}
            <Flex ml={4} _hover={{ cursor: "pointer" }}>
                {changeColorModeIcon}
            </Flex>
        </Flex>
    );
};

export default Header;

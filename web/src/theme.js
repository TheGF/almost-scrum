

import theme, { Theme } from "@chakra-ui/theme";
import color from "./color"


const styles = {
    ...theme.styles,
    global: (props) => ({
        ...theme.styles.global,
        fontFamily: "body",
        fontWeight: "light",
        "html, body": {
            background: color(props.colorMode, 'background'),
            color: color(props.colorMode, 'color'),
            lineHeight: "tall",
            height: '100%',
            display: 'flex',
            flexDirection: 'column',
        },
        ".react-tags": {
            background: color(props.colorMode, 'input'),
        },
        ".panel1": {
            background: color(props.colorMode, 'panel1bg'),
            marginTop: "4px",
            borderRadius: '5px',
        },
        ".panel2": {
            background: color(props.colorMode, 'panel2bg'),
            borderRadius: '2px',
        },
        ".te-editor": {
            background: color(props.colorMode, 'input'),
        },
        ".te-mode-switch-section": {
            background: color(props.colorMode, 'input'),
        },
        ".tui-editor-defaultUI-toolbar": {
            background: color(props.colorMode, 'panel2bg'),
        },
        ".react-tags__search-input": {
            backgroundColor: color(props.colorMode, 'input'),
            ".te-editor": {
                backgroundColor: color(props.colorMode, 'input'),
                color: color(props.colorMode, 'color'),
            },
            ".image-input": {
                backgroundColor: props.colorMode === "dark" ? "gray.400" : "gray.100",
            },
        },
    }),
};

const customTheme = {
    ...theme,
    fonts: {
        ...theme.fonts,
        body: `"Source Sans Pro",${theme.fonts.body}`,
        heading: `"Source Sans Pro",${theme.fonts.heading}`,
    },
    styles,
    components: {
        ...theme.components,
        Popover: {
            ...theme.components.Popover,
            variants: {
                responsive: {
                    popper: {
                        maxWidth: 'unset',
                        width: 'unset'
                    }
                }
            }
        }
    }
};

export default customTheme;

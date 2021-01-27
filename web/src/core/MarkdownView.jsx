
import { React } from "react";
import ReactMarkdown from "react-markdown";
import "react-mde/lib/styles/css/react-mde-all.css";
import gfm from 'remark-gfm';
import { Badge } from '@chakra-ui/react';


function MarkdownView(props) {
    const { highlights, ...more } = props

    function Image(props) {
        const token = localStorage.token
        if (token) {
            const src = `${props.src}?token=${token}`
            return <img {...props} style={{ maxWidth: '20%', maxHeight: '20%' }} src={src} />
        } else {
            return <img {...props} style={{ maxWidth: '20%', maxHeight: '20%' }} />
        }
    }

    function highlightWords(text) {
        let keyLocations = []
        const out = []
        const lower = text.toLowerCase()
        for (const key of highlights) {
            let start = 0
            let end = 0
            const lowerKey = key.toLowerCase()
            do {
                start = lower.indexOf(lowerKey, end)
                if (start >= 0) {
                    end = start + key.length
                    keyLocations.push([start, end])
                }
            } while (start != -1)
        }

        keyLocations = keyLocations.sort((a, b) => a[0] - b[0])
        let cursor = 0
        for (const [start, end] of keyLocations) {
            if (cursor > start) continue

            out.push(text.substring(cursor, start))
            out.push(<Badge key={out.length} colorScheme="green">
                {text.substring(start, end)}
            </Badge>)
            cursor = end
        }
        out.push(text.substring(cursor))
        return out
    }

    function TextFilter(props) {
        return highlights ? highlightWords(props.value) : props.value
    }

    const renderers = {
        image: Image,
        text: TextFilter,
    }


    return <ReactMarkdown
        plugins={[gfm]}
        renderers={renderers}
        {...more}
    />
}
export default MarkdownView
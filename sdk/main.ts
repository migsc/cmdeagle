import { argv } from "process";

type Params<A, F> = {
    args: Args<A>
    flags: Flags<F>
}

type Args<T> = T & {
    list: string[]
}

type Flags<T> = T 

const emptyParams: Params<any, any> = {args: {list:[]}, flags: {}}

export function readJSONParams<A, F = any>() {
    let input: Params<A, F> = emptyParams;

    for(let pos = 0; pos < argv.length; pos++ ) {
        try {
            input = JSON.parse(argv[pos]);
            return input;
        } catch(e) {
            continue
        }
    }

    return emptyParams;
}

export function readJSONArgs<A>() {
    const input = readJSONParams<A, any>();
    return input.args;
}

export function readJSONFlags<F>() {
    const input = readJSONParams<any, F>();
    return input.flags;
}

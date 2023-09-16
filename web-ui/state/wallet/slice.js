import { createSlice } from "@reduxjs/toolkit";
import connectToWallet from "./thunks/connectWallet";

const initialState = {
    connecting: false,
    connected: false,
    typeWallet: "",
    address: "",
    balance: "",
}

export const slice = createSlice({
    name: 'wallet',
    initialState,
    reducers: {
        disconnectWallet: (state, action) => {
            state.connected = false
            state.address = ""
            state.typeWallet = ""
        },
    },
    extraReducers(builder) {
        builder.addCase(connectToWallet.pending, (state) => {
            state.connecting = true
        })
        builder.addCase(connectToWallet.fulfilled, (state, action) => {
            state.connecting = false
            state.address = action.payload.address
            state.typeWallet = action.payload.typeWallet
            state.balance = action.payload.balance
            state.connected = action.payload.connected
        })
    }
})

export const {
    disconnectWallet,
} = slice.actions;
export default slice.reducer;
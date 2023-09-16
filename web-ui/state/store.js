import { configureStore } from '@reduxjs/toolkit'
import walletReducer from "./wallet/slice";

export const store = configureStore({
    reducer: {
        wallet: walletReducer,
    },
    middleware: (getDefaultMiddleware) =>
        getDefaultMiddleware({
            serializableCheck: false,
        }),
})
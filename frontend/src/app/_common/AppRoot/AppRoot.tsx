"use client";

import { ReduxProvider } from "@/redux/store";

import React, { PropsWithChildren } from "react";

const AppRoot = ({ children }: PropsWithChildren) => {
  return <ReduxProvider>{children}</ReduxProvider>;
};

export default AppRoot;

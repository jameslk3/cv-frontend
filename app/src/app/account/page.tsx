"use client";
import React from "react";

import { useRouter } from "next/navigation";

import UAuthForm from "@/components/UserLoginOrCreate";
import { useAuth } from "@/app/context/AuthContext";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";

import { Button } from "@/components/ui/button";

export default function Account() {
  const {
    setEmail,
    setPassword,
    isLoggedIn,
    authEmail,
    login,
    logout,
    sendVerificationEmail,
    loading,
  } = useAuth();
  const router = useRouter();

  const handleFormSubmit = async (
    typeSumbit: string,
    email: string,
    password: string
  ) => {
    setEmail(email);
    setPassword(password);

    // If we are creating an account, we must send an email to verify the account
    if (typeSumbit === "CREATE") {
      // Redirect to verify email page and send the verification email
      const success = await sendVerificationEmail(email, password);
      if (!success) {
        return;
      }

      router.push("/account/verify-email?email=" + email);
    } else {
      // If we are logging in, we must call the login function
      login(email, password, "LOGIN");
    }
  };

  return (
    <>
      {!isLoggedIn ? (
        <>
          <div className="flex items-center">
            <h1 className="text-lg font-semibold md:text-2xl">Your Account</h1>
          </div>
          <div className="flex flex-1 items-start justify-center rounded-lg border border-primary border-dashed shadow-sm">
            <div className="flex flex-col items-center gap-1">
              <div className="flex-col items-center">
                <UAuthForm handleFormSubmit={handleFormSubmit} />
                {loading && <Skeleton className="w-full max-w-md h-12" />}
              </div>
            </div>
          </div>
        </>
      ) : (
        <>
          <div className="flex items-center">
            <h1 className="text-lg font-semibold md:text-2xl">Your Account</h1>
          </div>
          <div className="flex flex-1 items-start rounded-lg border border-primary border-dashed shadow-sm">
            <div className="flex flex-col gap-1 w-full">
              <div className="flex-col">
                <div className="flex flex-col gap-1 w-full">
                  <Card className="w-full">
                    <CardHeader className="flex flex-row items-center gap-5 bg-muted/50">
                      <div className="flex flex-col">
                        <CardTitle className="flex text-lg">
                          Welcome Court Visionary!
                        </CardTitle>
                      </div>
                    </CardHeader>
                    <CardContent className="p-6 text-sm">
                      <div className="grid gap-3">
                        <div className="font-semibold">Email</div>
                        <ul className="grid gap-3">
                          <li className="flex items-center justify-between">
                            <span className="text-muted-foreground">
                              {authEmail}
                            </span>
                          </li>
                        </ul>
                        <Separator className="my-2" />
                        <div className="font-semibold">Password</div>
                        <ul className="grid gap-3">
                          <li className="flex items-center justify-between">
                            <span className="text-muted-foreground">
                              ********
                            </span>
                            <span></span>
                          </li>
                        </ul>
                      </div>
                    </CardContent>
                    <CardFooter className="flex flex-row items-center border-t bg-muted/50 px-6 py-3 gap-5">
                      <div className="flex flex-col">
                        <Button className="max-w-lg" onClick={() => logout()}>
                          Logout
                        </Button>
                      </div>
                      {/* <div className="flex flex-row">
                        <Button variant="destructive" className="max-w-lg" onClick={() => logout()}>
                          Delete Account
                        </Button>
                      </div> */}
                    </CardFooter>
                  </Card>
                </div>
              </div>
            </div>
          </div>
        </>
      )}
    </>
  );
}

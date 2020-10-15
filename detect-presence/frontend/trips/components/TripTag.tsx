import React from "react";

export function TripTag({tag}: { tag: string }) {
    return (
        <span
            className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium leading-4 bg-blue-100 text-blue-800">
            {tag}
        </span>
    );
}

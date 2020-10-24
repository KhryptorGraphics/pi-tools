import React from "react";
import useSWR from "swr";
import {Link} from "com_github_mjm_pi_tools/go-links/proto/links/links_pb";
import {LIST_RECENT_LINKS} from "com_github_mjm_pi_tools/homebase/go-links/lib/fetch";

export function RecentLinksList() {
    const {data, error} = useSWR<Link[]>(LIST_RECENT_LINKS);
    if (error) {
        console.error(error);
    }

    return (
        <div className="bg-white overflow-hidden shadow rounded-lg">
            <div className="px-4 py-5 sm:px-6">
                <h3 className="text-lg leading-6 font-medium text-gray-900">
                    Recently added links
                </h3>
                {data && data.map(trip => (
                    <div key={trip.getId()}>
                        {trip.getShortUrl()}
                    </div>
                ))}
            </div>

        </div>
    );
}

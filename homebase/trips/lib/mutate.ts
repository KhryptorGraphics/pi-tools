import {client} from "com_github_mjm_pi_tools/homebase/trips/lib/trips_client";
import {UpdateTripTagsRequest} from "com_github_mjm_pi_tools/detect-presence/proto/trips/trips_pb";

export async function updateTripTags(id: string, oldTags: readonly string[], newTags: string[]): Promise<void> {
    const oldTagsSet = new Set(oldTags);
    const newTagsSet = new Set(newTags);

    const tagsToAdd = newTags.filter(tag => !oldTagsSet.has(tag));
    const tagsToRemove = oldTags.filter(tag => !newTagsSet.has(tag));

    const req = new UpdateTripTagsRequest();
    req.setTripId(id);
    req.setTagsToAddList(tagsToAdd);
    req.setTagsToRemoveList(tagsToRemove);

    await client.updateTripTags(req);
}

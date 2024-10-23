import type { WorkflowRuns } from "./workflow-runs.ts";

// https://docs.github.com/en/rest/actions/workflow-runs?apiVersion=2022-11-28#list-workflow-runs-for-a-repository
const runs: WorkflowRuns = await fetch(
  `https://api.github.com/repos/${Deno.env.get("OWNER")}/${
    Deno.env.get("REPO")
  }/actions/runs`,
  {
    method: "GET",
    headers: {
      Accept: "application/vnd.github+json",
      Authorization: `Bearer ${Deno.env.get("GITHUB_TOKEN")}`,
      "X-Github-Api-Version": "2022-11-28",
    },
  },
).then((res) => res.json()).then((res) => res as WorkflowRuns);
console.log(`Total Count: ${runs.total_count}`);
for (const run of runs.workflow_runs) {
  // https://docs.github.com/en/rest/actions/workflow-runs?apiVersion=2022-11-28#delete-a-workflow-run
  const delResult: Response = await fetch(
    `https://api.github.com/repos/${Deno.env.get("OWNER")}/${
      Deno.env.get("REPO")
    }/actions/runs/${run.id}`,
    {
      method: "DELETE",
      headers: {
        Accept: "application/vnd.github+json",
        Authorization: `Bearer ${Deno.env.get("GITHUB_TOKEN")}`,
        "X-Github-Api-Version": "2022-11-28",
      },
    },
  );
  if (delResult.status == 204) {
    console.log(`204 - ${run.name} (${run.id})`);
  } else {
    console.log(
      `${delResult.status} - ${run.name} (${run.id}) - ${await delResult
        .text()}`,
    );
  }
}

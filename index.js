import "dotenv/config";
// https://docs.github.com/en/rest/actions/workflow-runs?apiVersion=2022-11-28#list-workflow-runs-for-a-repository
let runs = await fetch(
  `https://api.github.com/repos/${process.env.OWNER}/${process.env.REPO}/actions/runs`,
  {
    method: "GET",
    headers: {
      Accept: "application/vnd.github+json",
      Authorization: `Bearer ${process.env.GITHUB_TOKEN}`,
      "X-Github-Api-Version": "2022-11-28",
    },
  },
).then((res) => res.json());
console.log(`Total Count: ${runs.total_count}`);
for (let run of runs.workflow_runs) {
  // https://docs.github.com/en/rest/actions/workflow-runs?apiVersion=2022-11-28#delete-a-workflow-run
  let delResult = await fetch(
    `https://api.github.com/repos/${process.env.OWNER}/${process.env.REPO}/actions/runs/${run.id}`,
    {
      method: "DELETE",
      headers: {
        Accept: "application/vnd.github+json",
        Authorization: `Bearer ${process.env.GITHUB_TOKEN}`,
        "X-Github-Api-Version": "2022-11-28",
      },
    },
  );
  if (delResult.status == 204) {
    console.log(`204 - ${run.name} (${run.id})`);
  } else {
    console.log(
      `${delResult.status} - ${run.name} (${run.id}) - ${await delResult.text()}`,
    );
  }
}
